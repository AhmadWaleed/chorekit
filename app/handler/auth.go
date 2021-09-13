package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/gorilla/sessions"
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
	sess := ctx.SessionStore(c)

	type user struct {
		Name     string `form:"name" validate:"required"`
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	usr := new(user)
	if err := c.Bind(usr); err != nil {
		c.Logger().Error(err)
		sess.FlashError(http.StatusText(http.StatusBadRequest))
		return c.Render(http.StatusUnprocessableEntity, "auth/signin", map[string]string{
			"name":  usr.Name,
			"email": usr.Email,
		})
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}
		return c.Render(http.StatusUnprocessableEntity, "auth/signup", map[string]string{
			"name":  usr.Name,
			"email": usr.Email,
		})
	}

	hash, err := core.NewHasher().Generate(usr.Password)
	if err != nil {
		c.Logger().Error(err)
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
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
		sess.FlashError(errors.ErrorText(errors.DeplicateUserFound))
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
		sess.FlashError(errors.ErrorText(errors.EntityCreationError))
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
	sess := ctx.SessionStore(c)

	type user struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	usr := new(user)
	if err := c.Bind(usr); err != nil {
		sess.FlashError(http.StatusText(http.StatusBadRequest))
		return c.Render(http.StatusUnprocessableEntity, "auth/signin", map[string]string{
			"email": usr.Email,
		})
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}
		return c.Render(http.StatusUnprocessableEntity, "auth/signin", map[string]string{
			"email": usr.Email,
		})
	}

	store := ctx.Store(ctx.App.DB())
	dbuser := new(database.User)
	if err := store.User.First(dbuser); err != nil {
		sess.FlashError(errors.ErrorText(errors.UserNotFound))
		return c.Render(http.StatusUnprocessableEntity, "auth/signin", map[string]string{
			"email": usr.Email,
		})
	}

	if err := sess.Authenticate(*dbuser, func(s *sessions.Session) {
		s.Options = &sessions.Options{
			Domain:   "localhost",
			Path:     "/",
			MaxAge:   3600 * 8,
			HttpOnly: true,
			Secure:   false,
		}
	}); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/home")
}
