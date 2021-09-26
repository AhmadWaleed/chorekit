package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ahmadwaleed/chore/pkg/config"
	"github.com/ahmadwaleed/chore/pkg/executer"
	"github.com/ahmadwaleed/chore/pkg/ssh"
	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/labstack/echo/v4"
)

func RunGet(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	id, _ := strconv.Atoi(c.Param("id"))
	run := new(database.Run)
	if err := store.Run.First(run, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "/run/show", run)
}

func RunPost(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())
	sess := ctx.SessionStore(c)

	id, _ := strconv.Atoi(c.Param("id"))

	task := new(database.Task)
	if err := store.Task.First(task, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	go run(*task, c, store)

	sess.FlashSuccess("Task ran successfully.")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/show/%d", task.ID))
}

func run(t database.Task, c echo.Context, store *database.Store) {
	// create ssh servers config with temporary private key files for task servers.
	var hosts []ssh.Config
	var privKeys []*os.File
	for _, srv := range t.Servers {
		f, err := ioutil.TempFile("", "id_rda_")
		if err != nil {
			c.Logger().Error(err)
		}
		f.WriteString(srv.SSHPrivateKey)
		privKeys = append(privKeys, f)
		hosts = append(hosts, ssh.Config{
			User:   srv.User,
			Host:   srv.IP,
			Port:   strconv.Itoa(srv.Port),
			RSA_ID: f.Name(),
		})
	}
	// remove the temporary private key files after running tasks
	defer func(files []*os.File) {
		for _, f := range files {
			if err := os.Remove(f.Name()); err != nil {
				c.Logger().Error(err)
			}
		}
	}(privKeys)

	task := config.Task{
		Name:     t.Name,
		Env:      config.EnvVar(t.EnvVar()),
		Commands: strings.Split(t.Script, "\n"),
		Hosts:    hosts,
	}
	runner := executer.New("ssh")
	err := runner.Run(task, func(o *executer.CmdOutput) {
		var stdout, stderr string
		if o.Stdout.Len() > 0 {
			stdout = fmt.Sprintf("[%s](%s):\n%s", o.Host, o.Command, o.Stdout.String())
		}
		if o.Stderr.Len() > 0 {
			stderr = fmt.Sprintf("[%s](%s):\n%s", o.Host, o.Command, o.Stderr.String())
		}

		if err := store.Run.Create(&database.Run{
			TaskID: t.ID,
			Output: fmt.Sprintf("%s\n%s", stdout, stderr),
		}); err != nil {
			c.Logger().Error(err)
		}
	})
	if err != nil {
		c.Logger().Error(err)
	}
}
