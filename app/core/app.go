package core

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/ahmadwaleed/choreui/app/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Application struct {
	Echo          *echo.Echo        // HTTP middleware
	config        *config.AppConfig // Configuration
	db            *gorm.DB          // Database connection
	modelRegistry *models.Model     // Model registry for migration
}

// NewApp will create a new instance of the application
func BootstrapApp(config *config.AppConfig) *Application {
	app := &Application{}
	app.config = config
	i18n.Configure(config.LocaleDir, config.Lang, config.LangDomain)

	app.modelRegistry = models.NewModel()
	err := app.modelRegistry.OpenWithConfig(config)
	if err != nil {
		log.Fatalf("gorm: could not connect to db %q", err)
	}

	app.Echo = NewRouter(app)
	app.db = app.modelRegistry.DB

	return app
}

// GetModelRegistry returns the model registry
func (app *Application) ModelRegistry() *models.Model {
	return app.modelRegistry
}

// DB returns gorm (ORM)
func (app *Application) DB() *gorm.DB {
	return app.db
}

// DBConn returns gorm (ORM) underlyig sql deriver connection
func (app *Application) DBConn() (*sql.DB, error) {
	conn, err := app.db.DB()
	if err != nil {
		return nil, err
	}

	return conn, nil
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
		db, err := app.db.DB()
		if err != nil {
			app.Echo.Logger.Fatal(err)
		}

		if err := db.Close(); err != nil {
			app.Echo.Logger.Fatal(err)
		}
	}

	// shutdown http app
	if err := app.Echo.Shutdown(ctx); err != nil {
		app.Echo.Logger.Fatal(err)
	}
}
