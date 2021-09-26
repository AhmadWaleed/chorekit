package database

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Task struct {
	DBModel
	Servers []Server `gorm:"many2many:task_servers;"`
	Runs    []Run
	Name    string `sql:"type:varchar(30)"`
	Env     string
	Script  string
}

func (t *Task) EnvVar() map[string]string {
	slice := strings.Split(t.Env, ";")
	vars := make(map[string]string, len(slice))
	for _, e := range vars {
		env := strings.Split(e, "=")
		vars[env[0]] = env[1]
	}

	return vars
}

type TaskServer struct{ DBModel }

type TaskModel interface {
	First(m *Task, conds ...interface{}) error
	Find(m *[]Task) error
	FindMany(m *[]Task, ids []uint) error
	Create(m *Task) error
	Update(m *Task) error
	Delete(m *Task, conds ...interface{}) error
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

func (m *MysqlTaskModel) FindMany(t *[]Task, ids []uint) error {
	return m.DB.Where("ID IN ?", ids).Find(t).Error
}

func (m *MysqlTaskModel) Update(t *Task) error {
	return m.DB.Save(t).Error
}

func (m *MysqlTaskModel) Delete(t *Task, conds ...interface{}) error {
	return m.DB.Select("Runs", "Servers").Delete(t, conds...).Error
}
