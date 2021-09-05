package database

import "gorm.io/gorm"

type Store struct {
	User UserModel
	Host HostModel
}

type StoreFunc func(db *gorm.DB) *Store

func NewStoreFunc(db *gorm.DB) *Store {
	return &Store{
		User: &MysqlUserModel{db},
		Host: &MysqlHostModel{db},
	}
}
