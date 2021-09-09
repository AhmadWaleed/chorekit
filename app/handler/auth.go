package handler

import (
	"fmt"
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
	return c.Render(http.StatusOK, "auth.auth/signup", AuthViewModel{})
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
		c.Logger().Error(err)

		return c.Render(http.StatusUnprocessableEntity, "auth.auth/login", AuthViewModel{
			User:   User{Email: usr.Email},
			Errors: []string{http.StatusText(http.StatusBadRequest)},
		})
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)

		return c.Render(http.StatusUnprocessableEntity, "auth.auth/signup", AuthViewModel{
			User:   User{Name: usr.Name, Email: usr.Email},
			Errors: errs,
		})
	}

	hash, err := core.NewHasher().Generate(usr.Password)
	if err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusUnprocessableEntity, "auth.auth/signup", AuthViewModel{
			User:   User{Name: usr.Name, Email: usr.Email},
			Errors: []string{fmt.Sprintf("%s: %v", errors.ErrorText(errors.EntityCreationError), err)},
		})
	}

	store := ctx.Store(ctx.App.DB())
	err = store.User.Create(&database.User{
		Name:     usr.Name,
		Email:    usr.Email,
		Password: hash,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusUnprocessableEntity, "auth.auth/signup", AuthViewModel{
			User:   User{Name: usr.Name, Email: usr.Email},
			Errors: []string{errors.ErrorText(errors.EntityCreationError)},
		})
	}

	return c.Redirect(http.StatusSeeOther, "/home")
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
		return c.Render(http.StatusUnprocessableEntity, "auth.auth/login", AuthViewModel{
			User:   User{Email: usr.Email},
			Errors: []string{http.StatusText(http.StatusBadRequest)},
		})
	}

	if errs := ctx.App.Validator.Validate(usr); len(errs) > 0 {
		c.Logger().Error(errs)

		return c.Render(http.StatusUnprocessableEntity, "auth.auth/login", AuthViewModel{
			User:   User{Email: usr.Email},
			Errors: errs,
		})
	}

	store := ctx.Store(ctx.App.DB())
	dbuser := new(database.User)
	if err := store.User.First(dbuser); err != nil {
		return c.Render(http.StatusUnprocessableEntity, "auth.auth/login", AuthViewModel{
			User:   User{Email: usr.Email},
			Errors: []string{errors.ErrorText(errors.UserNotFound)},
		})
	}

	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusInternalServerError, "auth.auth/login", AuthViewModel{
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
	sess.Values["user"] = User{Name: dbuser.Name, Email: dbuser.Email}
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/home")
}
