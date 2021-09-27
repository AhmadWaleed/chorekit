package core

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"time"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/core/validator"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/labstack/echo/v4"
)

type Application struct {
	Echo      *echo.Echo // HTTP middleware
	Validator *validator.Validator
	config    *config.AppConfig // Configuration
	db        *sql.DB           // Database connection
}

// NewApp will create a new instance of the application
func NewApp(config *config.AppConfig) (*Application, error) {
	app := &Application{}
	app.config = config
	i18n.Configure(config.LocaleDir, config.Lang, config.LangDomain)

	db, err := database.Connect(config.ConnectionString)
	if err != nil {
		return nil, err
	}
	app.db = db

	v, err := validator.NewValidator()
	if err != nil {
		return nil, err
	}
	app.Validator = v
	app.Echo = NewRouter(app)

	return app, nil
}

// DB returns gorm (ORM)
func (app *Application) DB() *sql.DB {
	return app.db
}

// Config return the current app configuration
func (app *Application) Config() *config.AppConfig {
	return app.config
}

// Start the http app
func (app *Application) Start(addr string) error {
	return app.Echo.Start(addr)
}

// ServeStaticFiles serve static files for development purpose
func (app *Application) ServeStaticFiles() {
	app.Echo.Static("/assets", app.config.AssetsBuildDir)
}

// GracefulShutdown Wait for interrupt signal
// to gracefully shutdown the app with a timeout of 5 seconds.
func (app *Application) GracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close database connection
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			app.Echo.Logger.Fatal(err)
		}
	}

	// shutdown http app
	if err := app.Echo.Shutdown(ctx); err != nil {
		app.Echo.Logger.Fatal(err)
	}
}
