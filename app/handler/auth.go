package handler

import (
	"fmt"
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/core/view"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type User struct {
	Name  string
	Email string
}

func SignupGet(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/signup", nil)
}

func SignupPost(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess, _ := session.Get("session", ctx)

	type user struct {
		Name     string `form:"name" validate:"required"`
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	usr := new(user)
	if err := c.Bind(usr); err != nil {
		c.Logger().Error(err)
		sess.AddFlash(view.Flash{
			Message: http.StatusText(http.StatusBadRequest),
			Type:    view.FlashError,
		})
		sess.Save(c.Request(), c.Response())
		return c.Render(http.StatusUnprocessableEntity, "auth/login", map[string]string{
			"name":  usr.Name,
			"email": usr.Email,
		})
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)

		for _, err := range errs {
			sess.AddFlash(view.Flash{
				Message: err,
				Type:    view.FlashError,
			})
		}
		sess.Save(c.Request(), c.Response())
		return c.Render(http.StatusUnprocessableEntity, "auth/signup", map[string]string{
			"name":  usr.Name,
			"email": usr.Email,
		})
	}

	hash, err := core.NewHasher().Generate(usr.Password)
	if err != nil {
		c.Logger().Error(err)
		sess.AddFlash(view.Flash{
			Message: fmt.Sprintf("%s: %v", errors.ErrorText(errors.EntityCreationError), err),
			Type:    view.FlashError,
		})
		sess.Save(c.Request(), c.Response())
		return c.Render(http.StatusUnprocessableEntity, "auth/signup", map[string]string{
			"name":  usr.Name,
			"email": usr.Email,
		})
	}

	store := ctx.Store(ctx.App.DB())
	dbuser := new(database.User)
	err = store.User.First(dbuser, "email = ?", usr.Email)
	if err != nil {
		c.Logger().Error(err)
	}

	if err == nil {
		sess.AddFlash(view.Flash{
			Message: errors.ErrorText(errors.DeplicateUserFound),
			Type:    view.FlashError,
		})
		sess.Save(c.Request(), c.Response())
		return c.Render(http.StatusOK, "auth/signup", map[string]string{
			"name":  usr.Name,
			"email": usr.Email,
		})
	}

	err = store.User.Create(&database.User{
		Name:     usr.Name,
		Email:    usr.Email,
		Password: hash,
	})

	if err != nil {
		c.Logger().Error(err)
		sess.AddFlash(view.Flash{
			Message: errors.ErrorText(errors.EntityCreationError),
			Type:    view.FlashError,
		})
		sess.Save(c.Request(), c.Response())
		return c.Render(http.StatusUnprocessableEntity, "auth/signup", map[string]string{
			"name":  usr.Name,
			"email": usr.Email,
		})
	}

	return c.Redirect(http.StatusSeeOther, "/auth/signin")
}

func SignInGet(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/signin", nil)
}

func SignInPost(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess, _ := session.Get("session", ctx)

	type user struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	usr := new(user)
	if err := c.Bind(usr); err != nil {
		sess.AddFlash(view.Flash{
			Message: http.StatusText(http.StatusBadRequest),
			Type:    view.FlashError,
		})
		sess.Save(c.Request(), c.Response())

		return c.Render(http.StatusUnprocessableEntity, "auth/login", map[string]string{
			"email": usr.Email,
		})
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.AddFlash(view.Flash{
				Message: err,
				Type:    view.FlashError,
			})
		}
		sess.Save(c.Request(), c.Response())
		return c.Render(http.StatusUnprocessableEntity, "auth/login", map[string]string{
			"email": usr.Email,
		})
	}

	store := ctx.Store(ctx.App.DB())
	dbuser := new(database.User)
	if err := store.User.First(dbuser); err != nil {
		sess.AddFlash(view.Flash{
			Message: errors.ErrorText(errors.UserNotFound),
			Type:    view.FlashError,
		})
		sess.Save(c.Request(), c.Response())

		return c.Render(http.StatusUnprocessableEntity, "auth/login", map[string]string{
			"email": usr.Email,
		})
	}

	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error(err)
		sess.AddFlash(view.Flash{
			Message: http.StatusText(http.StatusInternalServerError),
			Type:    view.FlashError,
		})
		sess.Save(c.Request(), c.Response())

		return c.Render(http.StatusInternalServerError, "auth/login", map[string]string{
			"email": usr.Email,
		})
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["auth"] = true
	sess.Values["user"] = dbuser
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/home")
}
