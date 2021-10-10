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
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
)

var (
	srv struct {
		config *config.AppConfig
		logger *log.Logger
		app    *core.Application
	}
	pool     *dockertest.Pool
	resource *dockertest.Resource
	db       *sql.DB
	migrator *migrate.Migrate
)

func TestMain(m *testing.M) {
	srv.config = &config.AppConfig{
		ConnectionString:  "root:root@tcp(127.0.0.1:3306)/choreui_testing?charset=utf8mb4&parseTime=True&loc=Local",
		TemplateFolder:    "../../web/templates",
		TemplateLayoutDir: "layouts",
		TemplateExt:       "tmpl",
	}

	var err error
	srv.app, err = core.NewApp(srv.config)
	if err != nil {
		log.Fatalf("could not create new app: %v", err)
	}

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() error {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("could not connect to docker: %v", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		return fmt.Errorf("could not start resource: %v", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to docker: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	migrator, err = migrate.NewWithDatabaseInstance("file:///migrations", "mysql", driver)
	if err != nil {
		return fmt.Errorf("could not create migrator instance: %v", err)
	}

	if err := migrator.Up(); err != nil {
		return fmt.Errorf("could not run test db migrations: %v", err)
	}

	return nil
}

func teardown() error {
	if err := migrator.Drop(); err != nil {
		return fmt.Errorf("could not drop test db migrations: %v", err)
	}

	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("could not purge resource: %v", err)
	}

	return nil
}
