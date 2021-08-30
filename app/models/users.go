package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string `sql:"type:varchar(30)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type UserModel interface {
	First(m *User) error
	Find(m *[]User) error
	Create(m *User) error
	Ping() error
}

func NewUserStore(db *gorm.DB) UserModel {
	return &MysqlUserModel{db}
}

// UserStore implements the UserStore interface
type MysqlUserModel struct {
	DB *gorm.DB
}

func (s *MysqlUserModel) First(m *User) error {
	return s.DB.First(m).Error
}

func (s *MysqlUserModel) Create(m *User) error {
	return s.DB.Create(m).Error
}

func (s *MysqlUserModel) Find(m *[]User) error {
	return s.DB.Find(m).Error
}

func (s *MysqlUserModel) Ping() error {
	db, err := s.DB.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}
