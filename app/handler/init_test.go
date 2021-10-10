package handler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/core"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
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

	var err error
	e.app, err = core.NewApp(e.config)
	if err != nil {
		log.Fatalf("could not create new app: %v", err)
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
