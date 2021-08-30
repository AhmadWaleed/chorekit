package handler

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/labstack/echo/v4"
)

type healthcheckReport struct {
	Health  string          `json:"health"`
	Details map[string]bool `json:"details"`
}

// GetHealthcheck returns the current functional state of the application
func GetHealthcheck(c echo.Context) error {
	cc := c.(*core.AppContext)
	_ = cc
	m := healthcheckReport{Health: "OK"}

	// dbCheck := cc.UserStore.Ping()
	// cacheCheck := cc.Cache.Ping()

	// if dbCheck != nil {
	// 	m.Health = "NOT"
	// 	m.Details["db"] = false
	// }

	// if cacheCheck != nil {
	// 	m.Health = "NOT"
	// 	m.Details["cache"] = false
	// }

	return c.JSON(http.StatusOK, m)
}
