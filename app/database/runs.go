package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Run struct {
	gorm.Model
	TaskID uint
	Output string
}

type RunModel interface {
	First(m *Run, conds ...interface{}) error
	Find(m *[]Run) error
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

func (m *MysqlRunModel) Find(r *[]Run) error {
	return m.DB.Preload(clause.Associations).Find(r).Error
}

func (m *MysqlRunModel) Update(r *Run) error {
	return m.DB.Save(r).Error
}
