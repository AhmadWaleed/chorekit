package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Bucket struct {
	BaseModel
	Name  string
	Tasks []BucketTask
}

type BucketTask struct {
	gorm.Model
	Task     Task
	TaskID   uint
	BucketID uint
}

type BucketModel interface {
	First(b *Bucket, conds ...interface{}) error
	Find(b *[]Bucket) error
	FindMany(b *[]Bucket, ids []uint) error
	Create(b *Bucket) error
	Delete(m *Bucket, conds ...interface{}) error
	Update(b *Bucket) error
}

type MysqlBucketModel struct {
	DB *gorm.DB
}

func (b *MysqlBucketModel) First(buc *Bucket, conds ...interface{}) error {
	return b.DB.Preload(clause.Associations).First(buc, conds...).Error
}

func (b *MysqlBucketModel) Create(buc *Bucket) error {
	return b.DB.Create(buc).Error
}

func (b *MysqlBucketModel) Find(buc *[]Bucket) error {
	return b.DB.Preload("Tasks.Task").Preload(clause.Associations).Find(buc).Error
}

func (b *MysqlBucketModel) FindMany(buc *[]Bucket, ids []uint) error {
	return b.DB.Where("ID IN ?", ids).Find(buc).Error
}

func (b *MysqlBucketModel) Update(buc *Bucket) error {
	return b.DB.Save(buc).Error
}

func (m *MysqlBucketModel) Delete(b *Bucket, conds ...interface{}) error {
	if err := m.DB.Association("Tasks").Delete(b.Tasks); err != nil {
		return err
	}

	return m.DB.Delete(b, conds...).Error
}
