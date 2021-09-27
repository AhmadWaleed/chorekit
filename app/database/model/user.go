package model

import (
	"database/sql"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserModel interface {
	Create(name, email, password string) error
	FindByEmail(email string) (*User, error)
}

type MysqlUserModel struct{ *sql.DB }

func (m *MysqlUserModel) Create(name, email, password string) error {
	stmt, err := m.DB.Prepare(`INSERT INTO users(name, email, password) VALUES(?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, email, password)
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if int(driverErr.Number) == mysqlerr.ER_DUP_ENTRY {
			return ErrDuplicateEntity
		}
	}

	return err
}

func (m *MysqlUserModel) FindByEmail(email string) (*User, error) {
	u := new(User)
	row := m.DB.QueryRow(`SELECT * FROM users WHERE email = ?`, email)
	if err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}

		return nil, err
	}

	return u, row.Err()
}
