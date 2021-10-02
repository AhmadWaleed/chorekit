package database

import (
	"database/sql"

	"github.com/ahmadwaleed/choreui/app/database/model"
)

type Store struct {
	User   model.UserModel
	Server model.ServerModel
	Task   model.TaskModel
	Run    model.RunModel
	Bucket model.BucketModel
}

type StoreFunc func(db *sql.DB) *Store

func NewStoreFunc(db *sql.DB) *Store {
	return &Store{
		User:   &model.MysqlUserModel{DB: db},
		Server: &model.MysqlServerModel{DB: db},
		Task:   &model.MysqlTaskModel{DB: db},
		Run:    &model.MysqlRunModel{DB: db},
		Bucket: &model.MysqlBucketModel{DB: db},
	}
}
