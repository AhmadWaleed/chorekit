package database

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Task struct {
	ID        uint
	Servers   []Server `gorm:"many2many:task_servers;"`
	Name      string   `sql:"type:varchar(30)"`
	Env       string
	Script    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TaskServer struct {
	gorm.Model
}

type TaskModel interface {
	First(m *Task, conds ...interface{}) error
	Find(m *[]Task) error
	Create(m *Task) error
	Update(m *Task) error
	Ping() error
}

type MysqlTaskModel struct {
	DB *gorm.DB
}

func (m *MysqlTaskModel) First(t *Task, conds ...interface{}) error {
	return m.DB.Preload(clause.Associations).First(t, conds...).Error
}

func (m *MysqlTaskModel) Create(t *Task) error {
	return m.DB.Create(t).Error
}

func (m *MysqlTaskModel) Find(t *[]Task) error {
	return m.DB.Preload(clause.Associations).Find(t).Error
}

func (m *MysqlTaskModel) FindMany(t *[]Task, ids []int) error {
	return m.DB.Where("ID IN ?", ids).Find(t).Error
}

func (m *MysqlTaskModel) Update(t *Task) error {
	return m.DB.Save(t).Error
}

func (m *MysqlTaskModel) Ping() error {
	db, err := m.DB.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}
