package core

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core/session"
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

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.NewSessionStore(c)
			if sess.GetBool("Auth") {
				return next(c)
			}

			return echo.ErrUnauthorized
		}
	}
}

func GuestMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.NewSessionStore(c)
			if !sess.GetBool("Auth") {
				return next(c)
			}

			return c.Redirect(http.StatusSeeOther, "/home")
		}
	}
}
