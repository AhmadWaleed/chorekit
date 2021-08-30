package core

import (
	"net/http"

	"github.com/ahmadwaleed/choreui/app/core/errors"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v "gopkg.in/go-playground/validator.v9"
)

func NewRouter(app *Application) *echo.Echo {
	config := app.config
	e := echo.New()
	e.Validator = &Validator{validator: v.New()}

	cc := AppContext{
		App: app,
		Loc: i18n.New(),
	}

	e.Use(AppCtxMiddleware(&cc))

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
	renderer := newTemplateRenderer(config.LayoutDir, config.TemplateDir)
	e.Renderer = renderer

	return e
}

func httpErrHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	code := http.StatusInternalServerError

	switch v := err.(type) {
	case *echo.HTTPError:
		err := c.JSON(v.Code, v)
		if err != nil {
			c.Logger().Error("error handler: json encoding", err)
		}
	default:
		e := errors.NewBoom(errors.InternalError, "Bad implementation", nil)
		err := c.JSON(code, e)
		if err != nil {
			c.Logger().Error("error handler: json encoding", err)
		}
	}
}
