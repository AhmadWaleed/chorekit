package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/labstack/echo/v4"
)

type Task struct {
	Name      string `form:"name" validate:"required"`
	Env       string `form:"env"`
	Script    string `form:"script" validate:"required"`
	ServerIDs []int  `form:"servers" validate:"required"`
}

func CreateTaskGet(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	var servers []database.Server
	if err := store.Server.Find(&servers); err != nil {
		c.Logger().Error(err)
		sess.FlashWarning("Please create a server before creating tasks.")
		return c.Render(http.StatusOK, "task/create", nil)
	}

	return c.Render(http.StatusOK, "task/create", servers)
}

func CreateTaskPost(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	task := new(Task)
	if err := c.Bind(task); err != nil {
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	var servers []database.Server
	if err := store.Server.FindMany(&servers, task.ServerIDs); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	if errs := ctx.App.Validator.Validate(task); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}

		return c.Redirect(http.StatusSeeOther, "/tasks/create")
	}

	if err := store.Task.Create(&database.Task{
		Name:    task.Name,
		Env:     task.Env,
		Servers: servers,
		Script:  task.Script,
	}); err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
		return c.Redirect(http.StatusSeeOther, "task/create")
	}

	sess.FlashSuccess("Task created successfully.")
	return c.Redirect(http.StatusSeeOther, "/tasks/create")
}

func TaskIndex(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	var tasks []database.Task
	if err := store.Task.Find(&tasks); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "task/index", tasks)
}

func ShowTask(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	id, _ := strconv.Atoi(c.Param("id"))

	task := new(database.Task)
	if err := store.Task.First(task, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "task/show", task)
}
