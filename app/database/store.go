package database

import "gorm.io/gorm"

type Store struct {
	User UserModel
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		User: &MysqlUserModel{db},
	}
}
