package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	models "github.com/ahmadwaleed/choreui/app/database"
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"
)

type (
	User struct {
		Name     string `form:"name" validate:"required"`
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	AuthViewModel struct {
		Errors []string
	}
)

func SignupGet(c echo.Context) error {
	return c.Render(http.StatusOK, "auth.signup", AuthViewModel{})
}

func SignupPost(c echo.Context) error {
	ctx := c.(*core.AppContext)

	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := ctx.Echo().Validator.Validate(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return c.JSON(http.StatusBadRequest, err)
		}

		var vErrs []string
		trans := core.NewTranslator()
		for _, err := range err.(validator.ValidationErrors) {
			vErrs = append(vErrs, err.Translate(trans))
		}

		if len(vErrs) > 0 {
			return c.Render(http.StatusUnprocessableEntity, "auth.signup", AuthViewModel{vErrs})
		}
	}

	hash, err := core.NewHasher().Generate(user.Password)
	if err != nil {
		b := errors.NewBoom(errors.EntityCreationError, errors.ErrorText(errors.EntityCreationError), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, b)
	}

	err = ctx.Store.User.Create(&models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hash,
	})
	if err != nil {
		b := errors.NewBoom(errors.UserNotFound, errors.ErrorText(errors.UserNotFound), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, b)
	}

	return c.Render(http.StatusOK, "base.home", nil)
}
