package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ahmadwaleed/chore/pkg/config"
	"github.com/ahmadwaleed/chore/pkg/executer"
	choressh "github.com/ahmadwaleed/chore/pkg/ssh"
	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	Name          string `form:"name" validate:"required"`
	IP            string `form:"ip" validate:"required"`
	User          string `form:"user" validate:"required"`
	Port          int    `form:"port" validate:"required"`
	SSHPublicKey  string
	SSHPrivateKey string
	Status        string
}

func CreateServerGet(c echo.Context) error {
	return c.Render(http.StatusOK, "server/create", nil)
}

func CreateServerPost(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(ctx)

	srv := new(Server)
	if err := c.Bind(srv); err != nil {
		c.Logger().Error(err)
		sess.FlashError(http.StatusText(http.StatusBadRequest))
		return c.Render(http.StatusUnprocessableEntity, "server/create", nil)
	}

	if errs := ctx.App.Validator.Validate(srv); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}
		return c.Render(http.StatusUnprocessableEntity, "server/create", nil)
	}

	privKey, pubKey, err := generatePrivPubKeyPair()
	if err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
		return c.Render(http.StatusUnprocessableEntity, "server/create", nil)
	}

	store := ctx.Store(ctx.App.DB())
	err = store.Server.Create(&database.Server{
		Name:          srv.Name,
		IP:            srv.IP,
		User:          srv.User,
		Port:          srv.Port,
		SSHPrivateKey: privKey,
		SSHPublicKey:  pubKey,
		Status:        string(database.Inactive),
	})
	if err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
		return c.Render(http.StatusUnprocessableEntity, "server/create", nil)
	}

	return c.Render(http.StatusOK, "server/create", nil)
}

func DeleteServer(c echo.Context) error {
	ctx := c.(*core.AppContext)

	id, _ := strconv.Atoi(c.Param("id"))
	store := ctx.Store(ctx.App.DB())
	if err := store.Server.Delete(&database.Server{}, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusSeeOther, "/servers/index")
}

func ShowServer(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	id, _ := strconv.Atoi(c.Param("id"))

	server := new(database.Server)
	if err := store.Server.First(server, id); err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.InternalError))
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "server/show", server)
}

func IndexServer(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	var servers []database.Server
	if err := store.Server.Find(&servers); err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.InternalError))
		return c.Render(http.StatusOK, "server/index", nil)
	}

	return c.Render(http.StatusOK, "server/index", servers)
}

func StatusCheck(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	id, _ := strconv.Atoi(c.Param("id"))

	server := new(database.Server)
	if err := store.Server.First(server, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	file, err := ioutil.TempFile("", "id_rda_")
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(server.SSHPrivateKey); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	h := choressh.Config{
		User:   server.User,
		Host:   server.IP,
		Port:   strconv.Itoa(server.Port),
		RSA_ID: file.Name(),
	}

	task := config.Task{
		Name:     "",
		Commands: []string{"echo 1"},
		Hosts:    []choressh.Config{h},
	}

	runner := executer.New("ssh")
	err = runner.Run(task, func(o *executer.CmdOutput) {
		if o.Stderr.String() != "" {
			server.Status = string(database.Inactive)
			if err := store.Server.Update(server); err != nil {
				c.Logger().Error(err)
			}
		} else {
			server.Status = string(database.Active)
			if err := store.Server.Update(server); err != nil {
				c.Logger().Error(err)
			}
		}
	})

	if err != nil {
		server.Status = string(database.Inactive)
		if err := store.Server.Update(server); err != nil {
			c.Logger().Error(err)
		}
	}

	return c.Redirect(http.StatusSeeOther, "/servers/index")
}

func generatePrivPubKeyPair() (string, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("could not generate RSA keypair: %v", err)
	}

	privKeyb, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", "", fmt.Errorf("could not marshal PKCS8 private key: %v", err)
	}

	privKeyStr := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyb,
	}))

	pub, err := ssh.NewPublicKey(key.Public())
	if err != nil {
		return "", "", fmt.Errorf("could not generate SSH public key: %v", err)
	}
	pubKeyStr := string(ssh.MarshalAuthorizedKey(pub))

	return privKeyStr, pubKeyStr, nil
}
