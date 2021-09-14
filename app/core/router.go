package core

import (
	"fmt"
	"net/http"

	sess "github.com/ahmadwaleed/choreui/app/core/session"
	"github.com/ahmadwaleed/choreui/app/core/view"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(app *Application) *echo.Echo {
	config := app.config
	e := echo.New()

	e.Use(AppCtxMiddleware(&AppContext{
		App:          app,
		Loc:          i18n.New(),
		Store:        database.NewStoreFunc,
		SessionStore: sess.NewSessionStore,
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(app.config.AppKey))))

	if config.RequestLogger {
		e.Use(middleware.Logger()) // request logger
	}

	e.Use(middleware.Recover())       // panic errors are thrown
	e.Use(middleware.BodyLimit("5M")) // limit body payload to 5MB
	e.Use(middleware.Secure())        // provide protection against injection attacks
	e.Use(middleware.RequestID())     // generate unique requestId

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// add custom error formating
	e.HTTPErrorHandler = httpErrHandler

	// Add html templates with go template syntax
	e.Renderer = view.NewTemplate(app.Config())

	return e
}

func httpErrHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	code := http.StatusInternalServerError

	switch v := err.(type) {
	case *echo.HTTPError:
		errpage := fmt.Sprintf("web/templates/errors/%d.html", v.Code)
		if err := c.File(errpage); err != nil {
			c.Logger().Error(err)
		}
		c.Logger().Error(err)
	default:
		errpage := fmt.Sprintf("web/templates/errors/%d.html", code)
		if err := c.File(errpage); err != nil {
			c.Logger().Error(err)
		}
		c.Logger().Error(err)
	}
}
