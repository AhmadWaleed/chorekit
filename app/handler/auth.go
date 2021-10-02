package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database/model"
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
		return c.Render(http.StatusUnprocessableEntity, "auth/signup", nil)
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}
		return c.Render(http.StatusUnprocessableEntity, "auth/signup", nil)
	}

	hash, err := core.NewHasher().Generate(usr.Password)
	if err != nil {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityCreation.Error())
		return c.Render(http.StatusUnprocessableEntity, "auth/signup", nil)
	}

	store := ctx.Store(ctx.App.DB())
	_, err = store.User.FindByEmail(usr.Email)
	if err != nil && err != model.ErrNoResult {
		c.Logger().Error(err)
		sess.FlashError(model.ErrEntityCreation.Error())
		return c.Redirect(http.StatusSeeOther, "/auth/signup")
	}

	err = store.User.Create(usr.Name, usr.Email, hash)
	if err != nil {
		c.Logger().Error(err)
		if err == model.ErrDuplicateEntity {
			sess.FlashError(model.ErrDuplicateEntity.Error())
		}

		return c.Redirect(http.StatusSeeOther, "/auth/signup")
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
		return c.Redirect(http.StatusSeeOther, "/auth/signin")
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)
		for _, err := range errs {
			sess.FlashError(err)
		}
		return c.Redirect(http.StatusSeeOther, "/auth/signin")
	}

	store := ctx.Store(ctx.App.DB())
	dbuser, err := store.User.FindByEmail(usr.Email)
	if err != nil {
		c.Logger().Error(err)
		if err == model.ErrNoResult {
			sess.FlashError(model.ErrNoResult.Error())
		} else {
			sess.FlashError(http.StatusText(http.StatusInternalServerError))
		}
		return c.Redirect(http.StatusSeeOther, "/auth/signin")
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

func Signout(c echo.Context) error {
	ctx := c.(*core.AppContext)
	sess := ctx.SessionStore(ctx)
	sess.Logout()

	return c.Redirect(http.StatusSeeOther, "/auth/signin")
}
