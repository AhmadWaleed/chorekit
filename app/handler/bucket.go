package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database/model"
	"github.com/ahmadwaleed/choreui/app/utils"
	"github.com/labstack/echo/v4"
)

func CreateBucketGet(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	tasks, err := store.Task.All()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "bucket/create", tasks)
}

type Bucket struct {
	Name     string `form:"name" validate:"required"`
	Parallel bool   `form:"parallel"`
	TaskIDs  []uint `form:"tasks" validate:"required"`
}

func CreateBucketPost(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	bucket := new(Bucket)
	if err := c.Bind(bucket); err != nil {
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	if errs := ctx.App.Validator.Validate(bucket); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.create.get"))
	}

	tasks, err := store.Task.FindMany(bucket.TaskIDs)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	task := &model.Bucket{
		Name:     bucket.Name,
		Parallel: bucket.Parallel,
	}

	for _, t := range tasks {
		task.Tasks = append(task.Tasks, &model.BucketTask{Task: t})
	}

	if err := store.Bucket.Create(task); err != nil {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityCreation.Error())
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.create.get"))
	}

	sess.FlashSuccess("Bucket created successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.index"))
}

func UpdateBucket(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	ID, _ := strconv.Atoi(c.Param("id"))

	b := new(Bucket)
	if err := c.Bind(b); err != nil {
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	if errs := ctx.App.Validator.Validate(b); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/buckets/show/%d", ID))
	}

	bucket, err := store.Bucket.FindByID(uint(ID))
	if err != nil {
		c.Logger().Error(err)
		if err == model.ErrNoResult {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	tasks, err := store.Task.FindMany(b.TaskIDs)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	bucket.Name = b.Name
	bucket.Parallel = b.Parallel

	for _, t := range tasks {
		bucket.Tasks = append(bucket.Tasks, &model.BucketTask{Task: t})
	}

	if err := store.Bucket.Update(bucket); err != nil {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityCreation.Error())
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.create.get"))
	}

	sess.FlashSuccess("Bucket created successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.index"))
}

type (
	BucketList     struct{ Items []BucketListItem }
	BucketListItem struct {
		Bucket model.Bucket
		Tasks  []string
	}
)

func NewBucketListViewModel(buckets []*model.Bucket) BucketList {
	var items []BucketListItem
	for _, b := range buckets {
		var names []string
		for _, t := range b.Tasks {
			names = append(names, t.Task.Name)
		}
		items = append(items, BucketListItem{Bucket: *b, Tasks: names})
	}

	return BucketList{Items: items}
}

func IndexBucket(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	buckets, err := store.Bucket.All()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "bucket/index", NewBucketListViewModel(buckets))
}

type (
	BucketTask struct {
		Task  model.Task
		Added bool
	}
	BucketViewModel struct {
		Bucket model.Bucket
		Tasks  []BucketTask
	}
)

func NewBucketViewModel(b model.Bucket, tasks []*model.Task) BucketViewModel {
	var isAdded = func(task model.Task) bool {
		for _, t := range b.Tasks {
			if t.Task.ID == task.ID {
				return true
			}
		}
		return false
	}

	vm := BucketViewModel{Bucket: b}
	for _, t := range tasks {
		vm.Tasks = append(vm.Tasks, BucketTask{
			Task:  *t,
			Added: isAdded(*t),
		})
	}

	return vm
}

func ShowBucket(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	ID, _ := strconv.Atoi(c.Param("id"))
	bucket, err := store.Bucket.FindByID(uint(ID))
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	tasks, err := store.Task.All()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "bucket/show", NewBucketViewModel(*bucket, tasks))
}

func DeleteBucket(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	ID, _ := strconv.Atoi(c.Param("id"))
	if err := store.Bucket.Delete(uint(ID)); err != nil {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityDeletion.Error())
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.index"))
	}

	sess.FlashSuccess("Bucket deleted successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.index"))
}
