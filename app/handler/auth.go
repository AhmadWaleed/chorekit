package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type (
	User struct {
		Name  string
		Email string
	}

	AuthViewModel struct {
		User   User
		Errors []string
	}
)

func SignupGet(c echo.Context) error {
	return c.Render(http.StatusOK, "auth.signup", AuthViewModel{})
}

func SignupPost(c echo.Context) error {
	ctx := c.(*core.AppContext)

	type user struct {
		Name     string `form:"name" validate:"required"`
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	usr := new(user)
	if err := c.Bind(usr); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := ctx.Echo().Validator.Validate(usr); err != nil {
		c.Logger().Error(err)
		errs := core.TransValidationErrors(err)

		if len(errs) > 0 {
			return c.Render(http.StatusUnprocessableEntity, "auth.signup", AuthViewModel{
				User:   User{Name: usr.Name, Email: usr.Email},
				Errors: errs,
			})
		}
	}

	hash, err := core.NewHasher().Generate(usr.Password)
	if err != nil {
		b := errors.NewBoom(errors.EntityCreationError, errors.ErrorText(errors.EntityCreationError), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, b)
	}

	err = ctx.Store.User.Create(&database.User{
		Name:     usr.Name,
		Email:    usr.Email,
		Password: hash,
	})
	if err != nil {
		b := errors.NewBoom(errors.UserNotFound, errors.ErrorText(errors.UserNotFound), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, b)
	}

	return c.Render(http.StatusOK, "base.home", nil)
}

func SignInGet(c echo.Context) error {
	return c.Render(http.StatusOK, "base.home", nil)
}

func SignInPost(c echo.Context) error {
	ctx := c.(*core.AppContext)

	type user struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	usr := new(user)
	if err := c.Bind(usr); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := ctx.Echo().Validator.Validate(usr); err != nil {
		c.Logger().Error(err)
		errs := core.TransValidationErrors(err)

		if len(errs) > 0 {
			return c.Render(http.StatusUnprocessableEntity, "auth.login", AuthViewModel{
				User:   User{Email: usr.Email},
				Errors: errs,
			})
		}
	}

	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusInternalServerError, "auth.login", AuthViewModel{
			User:   User{Email: usr.Email},
			Errors: []string{http.StatusText(http.StatusInternalServerError)},
		})
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["auth"] = true
	sess.Values["user"] = User{Email: usr.Email}
	sess.Save(c.Request(), c.Response())

	return nil
}
