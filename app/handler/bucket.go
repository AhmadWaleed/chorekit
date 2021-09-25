package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/utils"
	"github.com/labstack/echo/v4"
)

func CreateBucketGet(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	var tasks []database.Task
	if err := store.Task.Find(&tasks); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "bucket/create", tasks)
}

type Bucket struct {
	Name    string `form:"name" validate:"required"`
	TaskIDs []uint `form:"tasks" validate:"required"`
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

	var tasks []database.Task
	if err := store.Task.FindMany(&tasks, bucket.TaskIDs); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	var bTasks []database.BucketTask
	for _, task := range tasks {
		bTasks = append(bTasks, database.BucketTask{Task: task})
	}

	if err := store.Bucket.Create(&database.Bucket{
		Name:  bucket.Name,
		Tasks: bTasks,
	}); err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.create.get"))
	}

	sess.FlashSuccess("Bucket created successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.index"))
}

type (
	BucketList     struct{ Items []BucketListItem }
	BucketListItem struct {
		Bucket database.Bucket
		Tasks  []string
	}
)

func NewBucketListViewModel(buckets []database.Bucket) BucketList {
	var items []BucketListItem
	for _, b := range buckets {
		var names []string
		for _, t := range b.Tasks {
			names = append(names, t.Task.Name)
		}
		items = append(items, BucketListItem{Bucket: b, Tasks: names})
	}

	return BucketList{Items: items}
}

func IndexBucket(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.Store(ctx.App.DB())

	var buckets []database.Bucket
	if err := store.Bucket.Find(&buckets); err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "bucket/index", NewBucketListViewModel(buckets))
}

func DeleteBucket(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(c)
	store := ctx.Store(ctx.App.DB())

	id, _ := strconv.Atoi(c.Param("id"))
	bucket := new(database.Bucket)
	if err := store.Bucket.First(bucket, id); err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	if err := store.Bucket.Delete(bucket); err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityDeletionError))
		return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.index"))
	}

	sess.FlashSuccess("Bucket deleted successfully.")
	return c.Redirect(http.StatusSeeOther, utils.Route(c, "bucket.index"))
}
