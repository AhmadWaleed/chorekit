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
	db       *sql.DB
	migrator *migrate.Migrate
)

func TestMain(m *testing.M) {
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
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	srv.config = &config.AppConfig{
		TemplateFolder:    "../../web/templates",
		TemplateLayoutDir: "layouts",
		TemplateExt:       "tmpl",
		RequestLogger:     false,
		AppKey:            "test",
	}

	srv.app, err = core.NewAppWithDB(db, srv.config)
	if err != nil {
		log.Fatalf("could not create new app: %v", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	migrator, err = migrate.NewWithDatabaseInstance("file:///Users/ahmedwaleed/.go/src/github.com/choreui/db/migrations", "mysql", driver)
	if err != nil {
		log.Fatalf("could not create migrator instance: %v", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func migrateUp() error {
	if err := migrator.Up(); err != nil {
		return fmt.Errorf("could not run test db migrations: %v", err)
	}
	return nil
}

func migrateDown() error {
	if err := migrator.Down(); err != nil {
		return fmt.Errorf("could not run drop db migrations: %v", err)
	}
	return nil
}
