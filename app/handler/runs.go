package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ahmadwaleed/chore/pkg/executer"
	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database/model"
	"github.com/ahmadwaleed/choreui/app/ssh"
	"github.com/labstack/echo/v4"
)

func ShowTaskRun(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	ID, _ := strconv.Atoi(c.Param("id"))
	run, err := store.Run.FindByID(uint(ID))
	if err != nil {
		c.Logger().Error(err)
		if err == model.ErrNoResult {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "/run/show", run)
}

func RunBucket(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())
	sess := ctx.SessionStore(c)

	ID, _ := strconv.Atoi(c.Param("id"))
	bucket, err := store.Bucket.FindByID(uint(ID))
	if err != nil {
		c.Logger().Error(err)
		if err == model.ErrNoResult {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	for _, t := range bucket.Tasks {
		go func(b *model.Bucket, t *model.Task, c echo.Context) {
			db := ctx.Store(ctx.App.DB())
			// Fetch task from task model to preload task servers.
			task, err := db.Task.FindByID(t.ID)
			if err != nil {
				c.Logger().Error(err)
				return
			}

			runner := ssh.Runner{}
			runner.RunTask(*task, func(o *executer.CmdOutput) {
				var stdout, stderr string
				if o.Stdout.Len() > 0 {
					stdout = fmt.Sprintf("[%s](%s):\n%s", o.Host, o.Command, o.Stdout.String())
				}
				if o.Stderr.Len() > 0 {
					stderr = fmt.Sprintf("[%s](%s):\n%s", o.Host, o.Command, o.Stderr.String())
				}

				if err := db.Bucket.CreateRun(b, &model.Run{
					TaskID: t.ID,
					Output: fmt.Sprintf("%s\n%s", stdout, stderr),
				}); err != nil {
					c.Logger().Error(err)
				}
			})
			if err := runner.Close(); err != nil {
				c.Logger().Errorf("Could not close SSH runner: %v", err)
			}
		}(bucket, t.Task, c)
	}

	sess.FlashSuccess("Bucket ran successfully.")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/buckets/show/%d", bucket.ID))
}

func RunTask(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())
	sess := ctx.SessionStore(c)

	ID, _ := strconv.Atoi(c.Param("id"))
	task, err := store.Task.FindByID(uint(ID))
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	go func(t model.Task, c echo.Context) {
		db := ctx.Store(ctx.App.DB())
		runner := new(ssh.Runner)
		runner.RunTask(t, func(o *executer.CmdOutput) {
			var stdout, stderr string
			if o.Stdout.Len() > 0 {
				stdout = fmt.Sprintf("[%s](%s):\n%s", o.Host, o.Command, o.Stdout.String())
			}
			if o.Stderr.Len() > 0 {
				stderr = fmt.Sprintf("[%s](%s):\n%s", o.Host, o.Command, o.Stderr.String())
			}

			if err := db.Run.Create(&model.Run{
				TaskID: t.ID,
				Output: fmt.Sprintf("%s\n%s", stdout, stderr),
			}); err != nil {
				c.Logger().Error(err)
			}
		})
		if err := runner.Close(); err != nil {
			c.Logger().Errorf("Could not close SSH runner: %v", err)
		}
	}(*task, c)

	sess.FlashSuccess("Task ran successfully.")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/show/%d", task.ID))
}
