package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/labstack/echo/v4"
)

type (
	Host struct {
		Name          string `form:"name" validate:"required"`
		IP            string `form:"name" validate:"required"`
		User          string `form:"email" validate:"required"`
		Port          int    `form:"email" validate:"required"`
		SSHPublicKey  string
		SSHPrivateKey string
		Status        string
	}

	HostViewModel struct {
		Host   Host
		Errors []string
	}
)

func CreateHostGet(c echo.Context) error {
	return c.Render(http.StatusOK, "base.server/create", nil)
}

func CreateHostPost(c echo.Context) error {
	ctx := c.(*core.AppContext)
	_ = ctx

	h := new(Host)
	if err := c.Bind(h); err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusUnprocessableEntity, "base.create_host", HostViewModel{
			Host:   *h,
			Errors: []string{http.StatusText(http.StatusBadRequest)},
		})
	}

	if err := ctx.Echo().Validator.Validate(h); err != nil {
		c.Logger().Error(err)
		errs := core.TransValidationErrors(err)

		if len(errs) > 0 {
			return c.Render(http.StatusUnprocessableEntity, "auth.signup", HostViewModel{
				Host:   *h,
				Errors: errs,
			})
		}
	}

	store := ctx.Store(ctx.App.DB())
	err := store.Host.Create(&database.Host{
		Name:   h.Name,
		IP:     h.IP,
		User:   h.User,
		Port:   h.Port,
		Status: string(database.Inactive),
	})
	if err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusUnprocessableEntity, "auth.signup", HostViewModel{
			Host:   *h,
			Errors: []string{errors.ErrorText(errors.EntityCreationError)},
		})
	}

	return c.Render(http.StatusOK, "base.create_host", nil)
}
