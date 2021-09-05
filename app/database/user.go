package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint
	Name      string `sql:"type:varchar(30)"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserModel interface {
	First(m *User) error
	Find(m *[]User) error
	Create(m *User) error
	Ping() error
}

// UserStore implements the UserStore interface
type MysqlUserModel struct {
	DB *gorm.DB
}

func (m *MysqlUserModel) First(u *User) error {
	return m.DB.First(u).Error
}

func (m *MysqlUserModel) Create(u *User) error {
	return m.DB.Create(u).Error
}

func (m *MysqlUserModel) Find(u *[]User) error {
	return m.DB.Find(u).Error
}

func (m *MysqlUserModel) Ping() error {
	db, err := m.DB.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}
