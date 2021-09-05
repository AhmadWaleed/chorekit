package database

import (
	"time"

	"gorm.io/gorm"
)

type Status string

var (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type Host struct {
	ID            uint
	Name          string `sql:"type:varchar(30)"`
	IP            string
	User          string
	Port          int
	SSHPublicKey  string
	SSHPrivateKey string
	Status        string `sql:"type:ENUM('active','inactive')"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type HostModel interface {
	First(m *Host) error
	Find(m *[]Host) error
	Create(m *Host) error
	Update(m *Host) error
	Ping() error
}

type MysqlHostModel struct {
	DB *gorm.DB
}

func (m *MysqlHostModel) First(h *Host) error {
	return m.DB.First(h).Error
}

func (m *MysqlHostModel) Create(h *Host) error {
	return m.DB.Create(h).Error
}

func (m *MysqlHostModel) Find(h *[]Host) error {
	return m.DB.Find(h).Error
}

func (m *MysqlHostModel) Update(h *Host) error {
	return m.DB.Save(h).Error
}

func (m *MysqlHostModel) Ping() error {
	db, err := m.DB.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}
