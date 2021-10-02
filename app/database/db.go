package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Connect to the database
func Connect(DNS string) (*sql.DB, error) {
	// EX: DNS = "user:password@tcp(127.0.0.1:3306)/hello"
	db, err := sql.Open("mysql", DNS)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
