package handler

import (
	"fmt"
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	ctx := c.(*core.AppContext)
	store := ctx.SessionStore(c)
	fmt.Println(store.Session.Values)
	return c.Render(http.StatusOK, "home", nil)
}
