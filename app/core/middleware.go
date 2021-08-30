package core

import (
	"github.com/labstack/echo/v4"
)

func AppCtxMiddleware(cc *AppContext) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc.Context = c
			return h(cc)
		}
	}
}
