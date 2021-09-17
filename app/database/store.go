package database

import "gorm.io/gorm"

type Store struct {
	User   UserModel
	Server ServerModel
	Task   TaskModel
}

type StoreFunc func(db *gorm.DB) *Store

func NewStoreFunc(db *gorm.DB) *Store {
	return &Store{
		User:   &MysqlUserModel{db},
		Server: &MysqlServerModel{db},
		Task:   &MysqlTaskModel{db},
	}
}
