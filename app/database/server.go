package database

import (
	"gorm.io/gorm"
)

type Status string

var (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type Server struct {
	DBModel
	Name          string `sql:"type:varchar(30)"`
	IP            string
	User          string
	Port          int
	SSHPublicKey  string
	SSHPrivateKey string
	Status        string `sql:"type:ENUM('active','inactive')"`
}

type ServerModel interface {
	First(m *Server, conds ...interface{}) error
	Find(m *[]Server) error
	FindMany(m *[]Server, ids []int) error
	Create(m *Server) error
	Delete(m *Server, conds ...interface{}) error
	Update(m *Server) error
}

type MysqlServerModel struct {
	DB *gorm.DB
}

func (m *MysqlServerModel) First(s *Server, conds ...interface{}) error {
	return m.DB.First(s, conds...).Error
}

func (m *MysqlServerModel) Create(s *Server) error {
	return m.DB.Create(s).Error
}

func (m *MysqlServerModel) Find(s *[]Server) error {
	return m.DB.Find(s).Error
}

func (m *MysqlServerModel) FindMany(s *[]Server, ids []int) error {
	return m.DB.Where("ID IN ?", ids).Find(s).Error
}

func (m *MysqlServerModel) Update(s *Server) error {
	return m.DB.Save(s).Error
}

func (m *MysqlServerModel) Delete(s *Server, conds ...interface{}) error {
	return m.DB.Delete(s, conds...).Error
}
