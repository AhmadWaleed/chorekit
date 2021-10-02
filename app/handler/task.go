package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database/model"
	"github.com/ahmadwaleed/choreui/app/utils"
	"github.com/labstack/echo/v4"
)

type Task struct {
	Name      string `form:"name" validate:"required"`
	Env       string `form:"env"`
	Script    string `form:"script" validate:"required"`
	ServerIDs []uint `form:"servers" validate:"required"`
}

func CreateTaskGet(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	servers, err := store.Server.All()
	if err != nil {
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

	servers, err := store.Server.FindMany(task.ServerIDs)
	if err != nil {
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

	var cmds []string
	for _, c := range strings.Split(task.Script, "\n") {
		cmds = append(cmds, strings.TrimSpace(c))
	}

	if err := store.Task.Create(&model.Task{
		Name:    task.Name,
		Env:     task.Env,
		Servers: servers,
		Script:  strings.Join(cmds, "\n"),
	}); err != nil {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityCreation.Error())
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "task.create.get"))
	}

	sess.FlashSuccess("Task created successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "task.index"))
}

func TaskIndex(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	tasks, err := store.Task.All()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "task/index", tasks)
}

type (
	TaskServer struct {
		Server   model.Server
		Assigned bool
	}
	TaskViewModel struct {
		Task    model.Task
		Servers []TaskServer
	}
)

func NewTaskViewModel(task *model.Task, servers []*model.Server) TaskViewModel {
	var isAssigned = func(srv model.Server) bool {
		for _, s := range task.Servers {
			if s.ID == srv.ID {
				return true
			}
		}
		return false
	}

	vm := TaskViewModel{Task: *task}
	for _, s := range servers {
		vm.Servers = append(vm.Servers, TaskServer{
			Server:   *s,
			Assigned: isAssigned(*s),
		})
	}

	return vm
}

func ShowTask(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	ID, _ := strconv.Atoi(c.Param("id"))

	task, err := store.Task.FindByID(uint(ID))
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}
	task.Script = strings.TrimSpace(task.Script)

	servers, err := store.Server.All()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "task/show", NewTaskViewModel(task, servers))
}

func UpdateTask(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	ID, _ := strconv.Atoi(c.Param("id"))

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

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/show/%d", ID))
	}

	task, err := store.Task.FindByID(uint(ID))
	if err != nil {
		c.Logger().Error(err)
		if err == model.ErrNoResult {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	var cmds []string
	for _, c := range strings.Split(t.Script, "\n") {
		cmds = append(cmds, strings.TrimSpace(c))
	}

	servers, err := store.Server.FindMany(t.ServerIDs)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	task.Name = t.Name
	task.Env = t.Env
	task.Servers = servers
	task.Script = strings.Join(cmds, "\n")
	if err := store.Task.Update(task); err != nil {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityCreation.Error())
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/show/%d", task.ID))
	}

	sess.FlashSuccess("Task updated successfully.")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/show/%d", task.ID))
}

func DeleteTask(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	ID, _ := strconv.Atoi(c.Param("id"))

	if err := store.Task.Delete(uint(ID)); err != nil {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityCreation.Error())
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "task.index"))
	}

	sess.FlashSuccess("Task deleted successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "task.index"))
}
