package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/context"
	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/models"
	"github.com/labstack/echo/v4"
)

type (
	User          struct{}
	UserViewModel struct {
		Name string
		ID   string
	}
)

func GetUser(c echo.Context) error {
	cc := c.(*context.AppContext)
	userID := c.Param("id")

	user := models.User{ID: userID}

	err := cc.UserStore.First(&user)

	if err != nil {
		b := errors.NewBoom(errors.UserNotFound, errors.ErrorText(errors.UserNotFound), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, b)
	}

	vm := UserViewModel{
		Name: user.Name,
		ID:   user.ID,
	}

	return c.Render(http.StatusOK, "user.html", vm)

}

func GetUserJSON(c echo.Context) error {
	cc := c.(*context.AppContext)
	userID := c.Param("id")

	user := models.User{ID: userID}

	err := cc.UserStore.First(&user)

	if err != nil {
		b := errors.NewBoom(errors.UserNotFound, errors.ErrorText(errors.UserNotFound), err)
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, b)
	}

	vm := UserViewModel{
		Name: user.Name,
		ID:   user.ID,
	}

	return c.JSON(http.StatusOK, vm)
}
