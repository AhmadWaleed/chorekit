package database

import (
	"gorm.io/gorm"
)

type User struct {
	DBModel
	Name     string `sql:"type:varchar(30)"`
	Email    string `gorm:"unique"`
	Password string
}

type UserModel interface {
	First(m *User, conds ...interface{}) error
	Find(m *[]User) error
	Create(m *User) error
}

// UserStore implements the UserStore interface
type MysqlUserModel struct {
	DB *gorm.DB
}

func (m *MysqlUserModel) First(u *User, conds ...interface{}) error {
	return m.DB.First(u, conds...).Error
}

func (m *MysqlUserModel) Create(u *User) error {
	return m.DB.Create(u).Error
}

func (m *MysqlUserModel) Find(u *[]User) error {
	return m.DB.Find(u).Error
}
