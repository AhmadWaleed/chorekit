package main

import (
	"log"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/handler"
	"github.com/labstack/echo/v4"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// create server
	app, err := core.NewApp(config)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	// serve files for dev
	app.ServeStaticFiles()

	// register application routes
	RegisterRoutes(app.Echo)

	// metric / health endpoint according to RFC 5785
	app.Echo.GET("/.well-known/health-check", handler.GetHealthcheck)
	// app.Echo.GET("/.well-known/metrics", echo.WrapHandler(promhttp.Handler()))

	// Start server
	go func() {
		if err := app.Start(config.Address); err != nil {
			app.Echo.Logger.Info("shutting down the server")
		}
	}()

	app.GracefulShutdown()
}

func RegisterRoutes(e *echo.Echo) {
	// root routes
	e.GET("/home", handler.Home, core.AuthMiddleware()).Name = "home"
	e.GET("/signout", handler.Signout, core.AuthMiddleware()).Name = "logout"

	// auth endpoints
	auth := e.Group("/auth", core.GuestMiddleware())
	auth.GET("/signup", handler.SignupGet).Name = "signup.get"
	auth.POST("/signup", handler.SignupPost).Name = "signup.post"
	auth.GET("/signin", handler.SignInGet).Name = "sigin.get"
	auth.POST("/signin", handler.SignInPost).Name = "sigin.post"

	task := e.Group("/tasks")
	task.GET("/create", handler.CreateTaskGet).Name = "task.create.get"
	task.POST("/create", handler.CreateTaskPost).Name = "task.create.post"
	task.GET("/index", handler.TaskIndex).Name = "task.index"
	task.GET("/show/:id", handler.ShowTask).Name = "task.show"
	task.POST("/update/:id", handler.UpdateTask).Name = "task.update"
	task.POST("/delete/:id", handler.DeleteTask).Name = "task.delete"
	task.POST("/runs/:id", handler.RunPost).Name = "task.run"
	task.GET("/runs/show/:id", handler.RunGet).Name = "task.run.show"

	server := e.Group("/servers", core.AuthMiddleware())
	server.GET("/create", handler.CreateServerGet).Name = "server.create.get"
	server.POST("/create", handler.CreateServerPost).Name = "server.create.post"
	server.GET("/index", handler.IndexServer).Name = "server.index"
	server.GET("/show/:id", handler.ShowServer).Name = "server.show"
	server.POST("/delete/:id", handler.DeleteServer).Name = "server.delete"
	server.POST("/status/check/:id", handler.StatusCheck).Name = "server.status"

	bucket := e.Group("/buckets", core.AuthMiddleware())
	bucket.GET("/create", handler.CreateBucketGet).Name = "bucket.create.get"
	bucket.POST("/create", handler.CreateBucketPost).Name = "bucket.create.post"
	bucket.GET("/index", handler.IndexBucket).Name = "bucket.index"
	bucket.POST("/delete/:id", handler.DeleteBucket).Name = "bucket.delete"
	bucket.GET("/show/:id", handler.ShowBucket).Name = "bucket.show"
	bucket.POST("/update/:id", handler.UpdateBucket).Name = "bucket.update"
}
