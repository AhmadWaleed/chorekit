package handler

import (
	"log"
	"os"
	"testing"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/core"
	"github.com/ahmadwaleed/choreui/app/database"
)

var e struct {
	config *config.AppConfig
	logger *log.Logger
	app    *core.Application
}

func TestMain(m *testing.M) {
	e.config = &config.AppConfig{
		ConnectionString:  "root:root@tcp(127.0.0.1:3306)/choreui_testing?charset=utf8mb4&parseTime=True&loc=Local",
		TemplateFolder:    "../../web/templates",
		TemplateLayoutDir: "layouts",
		TemplateExt:       "tmpl",
	}

	e.app = core.NewApp(e.config)

	setup()
	c := m.Run()
	teardown()

	os.Exit(c)
}

func setup() {
	mr := e.app.ModelRegistry()
	if err := e.app.ModelRegistry().Register(database.User{}); err != nil {
		e.app.Echo.Logger.Fatal(err)
	}

	mr.AutoMigrate()
}

func teardown() {
	if err := e.app.ModelRegistry().AutoDropAll(); err != nil {
		e.app.Echo.Logger.Fatal(err)
	}
}
