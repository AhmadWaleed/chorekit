package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Name          string `form:"name" validate:"required"`
	IP            string `form:"name" validate:"required"`
	User          string `form:"email" validate:"required"`
	Port          int    `form:"email" validate:"required"`
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

	store := ctx.Store(ctx.App.DB())
	err := store.Server.Create(&database.Server{
		Name:   srv.Name,
		IP:     srv.IP,
		User:   srv.User,
		Port:   srv.Port,
		Status: string(database.Inactive),
	})
	if err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
		return c.Render(http.StatusUnprocessableEntity, "server/create", nil)
	}

	return c.Render(http.StatusOK, "server/create", nil)
}
