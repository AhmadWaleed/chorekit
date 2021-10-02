package utils

import (
	"github.com/labstack/echo/v4"
)

func Route(c echo.Context, name string) string {
	for _, r := range c.Echo().Routes() {
		if name == r.Name {
			return r.Path
		}
	}

	return ""
}
