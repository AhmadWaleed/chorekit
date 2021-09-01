package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/models"
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"
)

type user struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func Signup(c echo.Context) error {
	ctx := c.(*core.AppContext)
	_ = ctx

	user := new(user)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := ctx.Echo().Validator.Validate(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return c.JSON(http.StatusBadRequest, err)
		}

		var errors []string
		trans := core.NewTranslator()
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(trans))
		}

		if len(errors) > 0 {
			return c.JSON(http.StatusUnprocessableEntity, errors)
		}
	}

	store := models.NewUserStore(ctx.App.DB())
	if err := store.Create(&models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}); err != nil {
		b := errors.NewBoom(errors.UserNotFound, errors.ErrorText(errors.UserNotFound), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, b)
	}

	return nil
}
