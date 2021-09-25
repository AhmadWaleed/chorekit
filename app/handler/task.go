package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/utils"
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

		return c.Redirect(http.StatusSeeOther, utils.Route(c, "task.create.get"))
	}

	if err := store.Task.Create(&database.Task{
		Name:    task.Name,
		Env:     task.Env,
		Servers: servers,
		Script:  task.Script,
	}); err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "task.create.get"))
	}

	sess.FlashSuccess("Task created successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "task.index"))
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

type (
	TaskServer struct {
		Server   database.Server
		Assigned bool
	}
	TaskViewModel struct {
		Task    database.Task
		Servers []TaskServer
	}
)

func NewTaskViewModel(task database.Task, servers []database.Server) TaskViewModel {
	var isAssigned = func(srv database.Server) bool {
		for _, s := range task.Servers {
			if s.ID == srv.ID {
				return true
			}
		}
		return false
	}

	vm := TaskViewModel{Task: task}
	for _, s := range servers {
		vm.Servers = append(vm.Servers, TaskServer{
			Server:   s,
			Assigned: isAssigned(s),
		})
	}

	return vm
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
	task.Script = strings.TrimSpace(task.Script)

	var servers []database.Server
	if err := store.Server.Find(&servers); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "task/show", NewTaskViewModel(*task, servers))
}

func UpdateTask(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	t := new(Task)
	if err := c.Bind(t); err != nil {
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	if errs := ctx.App.Validator.Validate(t); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}

		return c.Redirect(http.StatusSeeOther, "/tasks/create")
	}

	id, _ := strconv.Atoi(c.Param("id"))
	task := new(database.Task)
	if err := store.Task.First(task, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	var cmds []string
	for _, c := range strings.Split(t.Script, "\n") {
		cmds = append(cmds, strings.TrimSpace(c))
	}

	var servers []database.Server
	if err := store.Server.FindMany(&servers, t.ServerIDs); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	task.Name = t.Name
	task.Env = t.Env
	task.Servers = servers
	task.Script = strings.Join(cmds, "\n")
	if err := store.Task.Update(task); err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
		return c.Redirect(http.StatusSeeOther, "task/create")
	}

	sess.FlashSuccess("Task updated successfully.")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/show/%d", task.ID))
}
