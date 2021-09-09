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

type Server struct {
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

type ServerModel interface {
	First(m *Server) error
	Find(m *[]Server) error
	Create(m *Server) error
	Update(m *Server) error
	Ping() error
}

type MysqlServerModel struct {
	DB *gorm.DB
}

func (m *MysqlServerModel) First(h *Server) error {
	return m.DB.First(h).Error
}

func (m *MysqlServerModel) Create(h *Server) error {
	return m.DB.Create(h).Error
}

func (m *MysqlServerModel) Find(h *[]Server) error {
	return m.DB.Find(h).Error
}

func (m *MysqlServerModel) Update(h *Server) error {
	return m.DB.Save(h).Error
}

func (m *MysqlServerModel) Ping() error {
	db, err := m.DB.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}
