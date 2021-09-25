package database

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Run struct {
	gorm.Model
	Task   Task
	TaskID uint
	Output string
}

func (r *Run) DisplayOutput() []string {
	return strings.Split(r.Output, "\n")
}

type RunModel interface {
	First(m *Run, conds ...interface{}) error
	Find(m *[]Run, conds ...interface{}) error
	Create(m *Run) error
	Update(m *Run) error
}

type MysqlRunModel struct {
	DB *gorm.DB
}

func (m *MysqlRunModel) First(r *Run, conds ...interface{}) error {
	return m.DB.Preload(clause.Associations).First(r, conds...).Error
}

func (m *MysqlRunModel) Create(r *Run) error {
	return m.DB.Create(r).Error
}

func (m *MysqlRunModel) Find(r *[]Run, conds ...interface{}) error {
	return m.DB.Preload(clause.Associations).Find(r, conds...).Error
}

func (m *MysqlRunModel) Update(r *Run) error {
	return m.DB.Save(r).Error
}
