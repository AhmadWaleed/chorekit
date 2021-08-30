package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
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
	cc := c.(*core.AppContext)
	userID := c.Param("id")

	user := models.User{ID: userID}

	store := models.NewUserStore(cc.App.DB())

	err := store.First(&user)
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
	cc := c.(*core.AppContext)
	userID := c.Param("id")

	user := models.User{ID: userID}

	store := models.NewUserStore(cc.App.DB())

	err := store.First(&user)
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
