package handler

import (
	"fmt"
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

func RunTask(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())
	id, _ := strconv.Atoi(c.Param("id"))

	t := new(database.Task)
	if err := store.Task.First(t, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	var hosts []ssh.Config
	for _, srv := range t.Servers {
		hosts = append(hosts, ssh.Config{
			User: srv.User,
			Host: srv.IP,
			Port: string(srv.Port),
		})
	}

	task := config.Task{
		Name:     t.Name,
		Env:      config.EnvVar(t.EnvVar()),
		Commands: strings.Split(t.Script, "\n"),
	}

	runner := executer.New("ssh")
	if err := runner.Run(task, func(o *executer.CmdOutput) { o.Display() }); err != nil {
		fmt.Fprintf(os.Stderr, "could not run task: %v\n", err)
	}

	return nil
}
