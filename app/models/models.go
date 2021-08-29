package models

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/ahmadwaleed/choreui/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Model facilitate database interactions
type Model struct {
	models map[string]reflect.Value
	isOpen bool
	*gorm.DB
}

// NewModel returns a new Model without opening database connection
func NewModel() *Model {
	return &Model{
		models: make(map[string]reflect.Value),
	}
}

// IsOpen returns true if the Model has already established connection
// to the database
func (m *Model) IsOpen() bool {
	return m.isOpen
}

// OpenWithConfig opens database connection with the settings found in cfg
func (m *Model) OpenWithConfig(cfg *config.AppConfig) error {
	conn, err := gorm.Open(mysql.Open(cfg.ConnectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	db, err := conn.DB()
	if err != nil {
		return err
	}

	// https://github.com/go-sql-driver/mysql/issues/461
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(20)

	m.DB = conn
	m.isOpen = true
	return nil
}

// Register adds the values to the models registry
func (m *Model) Register(values ...interface{}) error {

	// do not work on them.models first, this is like an insurance policy
	// whenever we encounter any error in the values nothing goes into the registry
	models := make(map[string]reflect.Value)
	if len(values) > 0 {
		for _, val := range values {
			rVal := reflect.ValueOf(val)
			if rVal.Kind() == reflect.Ptr {
				rVal = rVal.Elem()
			}
			switch rVal.Kind() {
			case reflect.Struct:
				models[getTypeName(rVal.Type())] = reflect.New(rVal.Type())
			default:
				return errors.New("models must be structs")
			}
		}
	}
	for k, v := range models {
		m.models[k] = v
	}
	return nil
}

// AutoMigrateAll runs migrations for all the registered models
func (m *Model) AutoMigrateAll() {
	for _, v := range m.models {
		m.AutoMigrate(v.Interface())
	}
}

// AutoDropAll drops all tables of all registered models
func (m *Model) AutoDropAll() {
	for _, v := range m.models {
		m.Migrator().DropTable(v)
	}
}

func getTypeName(typ reflect.Type) string {
	if typ.Name() != "" {
		return typ.Name()
	}
	split := strings.Split(typ.String(), ".")
	return split[len(split)-1]
}
