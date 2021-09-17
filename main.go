package main

import (
	"log"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/handler"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// create server
	app := core.NewApp(config)
	// serve files for dev
	app.ServeStaticFiles()

	// root routes
	app.Echo.GET("/home", handler.Home, core.AuthMiddleware())
	app.Echo.GET("/signout", handler.Signout, core.AuthMiddleware())

	// auth endpoints
	auth := app.Echo.Group("/auth", core.GuestMiddleware())
	auth.GET("/signup", handler.SignupGet)
	auth.POST("/signup", handler.SignupPost)
	auth.GET("/signin", handler.SignInGet)
	auth.POST("/signin", handler.SignInPost)

	task := app.Echo.Group("/tasks")
	task.GET("/create", handler.CreateTaskGet)
	task.POST("/create", handler.CreateTaskPost)
	task.GET("/index", handler.TaskIndex)
	task.GET("/show/:id", handler.ShowTask)

	host := app.Echo.Group("/servers", core.AuthMiddleware())
	host.GET("/create", handler.CreateServerGet)
	host.POST("/create", handler.CreateServerPost)
	host.GET("/index", handler.IndexServer)
	host.GET("/:id", handler.ShowServer)

	// api endpoints
	// g := app.Echo.Group("/api")
	// g.GET("/users/:id", handler.GetUserJSON)

	// // pages
	// u := app.Echo.Group("/users")
	// u.GET("/:id", handler.GetUser)

	// metric / health endpoint according to RFC 5785
	app.Echo.GET("/.well-known/health-check", handler.GetHealthcheck)
	// app.Echo.GET("/.well-known/metrics", echo.WrapHandler(promhttp.Handler()))

	// migration for dev
	mr := app.ModelRegistry()
	if err := mr.Register(database.User{}, database.Server{}, database.Task{}); err != nil {
		app.Echo.Logger.Fatal(err)
	}

	if err := mr.AutoDropAll(); err != nil {
		app.Echo.Logger.Fatal(err)
	}

	if err := mr.AutoMigrateAll(); err != nil {
		app.Echo.Logger.Fatal(err)
	}

	if err := mr.Register(database.User{}, database.Server{}, database.Task{}); err != nil {
		app.Echo.Logger.Fatal(err)
	}

	// Start server
	go func() {
		if err := app.Start(config.Address); err != nil {
			app.Echo.Logger.Info("shutting down the server")
		}
	}()

	app.GracefulShutdown()
}
