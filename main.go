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
	app.Echo.GET("/home", handler.Home)

	// auth endpoints
	auth := app.Echo.Group("/auth")
	auth.GET("/signup", handler.SignupGet)
	auth.POST("/signup", handler.SignupPost)
	auth.GET("/signin", handler.SignInGet)
	auth.POST("/signin", handler.SignInPost)

	host := app.Echo.Group("/servers")
	host.GET("/create", handler.CreateServerGet)
	host.POST("/create", handler.CreateServerPost)

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
	user := database.User{Name: "Peter"}
	mr := app.ModelRegistry()
	if err := mr.Register(user); err != nil {
		app.Echo.Logger.Fatal(err)
	}

	if err := mr.AutoDropAll(); err != nil {
		app.Echo.Logger.Fatal(err)
	}

	if err := mr.AutoMigrateAll(); err != nil {
		app.Echo.Logger.Fatal(err)
	}

	mr.Create(&user)

	// Start server
	go func() {
		if err := app.Start(config.Address); err != nil {
			app.Echo.Logger.Info("shutting down the server")
		}
	}()

	app.GracefulShutdown()
}
