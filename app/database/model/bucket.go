package model

import (
	"database/sql"
)

type Bucket struct {
	Model
	Name     string
	Parallel bool
	Tasks    []*BucketTask
}

type BucketTask struct {
	Model
	Task     *Task
	TaskID   uint
	BucketID uint
}

type BucketModel interface {
	Create(b *Bucket) error
	Update(b *Bucket) error
	Delete(ID uint) error
	All() ([]*Bucket, error)
	FindByID(ID uint) (*Bucket, error)
}

type MysqlBucketModel struct {
	DB *sql.DB
}

func (m *MysqlBucketModel) Create(b *Bucket) error {
	return trans(m.DB, func(tx *sql.Tx) error {
		result, err := tx.Exec(`INSERT INTO buckets (name, parallel) VALUES (?, ?)`, b.Name, b.Parallel)
		if err != nil {
			return err
		}
		ID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(`INSERT INTO bucket_tasks (task_id, bucket_id) VALUES (?, ?)`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		for _, task := range b.Tasks {
			_, err := stmt.Exec(task.Task.ID, ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *MysqlBucketModel) Update(b *Bucket) error {
	return trans(m.DB, func(tx *sql.Tx) error {
		_, err := tx.Exec(`Update buckets SET name=?, parallel=? WHERE id = ?`, b.Name, b.Parallel, b.ID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`DELETE FROM bucket_tasks WHERE bucket_id = ?`, b.ID)
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(`INSERT INTO bucket_tasks (task_id, bucket_id) VALUES(?, ?)`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		for _, task := range b.Tasks {
			_, err := stmt.Exec(task.Task.ID, b.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *MysqlBucketModel) All() ([]*Bucket, error) {
	rows, err := m.DB.Query(`SELECT * FROM buckets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buckets []*Bucket
	for rows.Next() {
		b := new(Bucket)
		if err := rows.Scan(
			&b.ID,
			&b.Name,
			&b.Parallel,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		buckets = append(buckets, b)
	}

	return buckets, rows.Err()
}

func (m *MysqlBucketModel) FindByID(ID uint) (*Bucket, error) {
	b := new(Bucket)
	row := m.DB.QueryRow(`SELECT * FROM buckets WHERE id = ?`, ID)
	if err := row.Scan(
		&b.ID,
		&b.Name,
		&b.Parallel,
		&b.CreatedAt,
		&b.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}

	return b, row.Err()
}

func (m *MysqlBucketModel) Delete(ID uint) error {
	return trans(m.DB, func(tx *sql.Tx) error {
		_, err := tx.Exec(`DELETE FROM buckets WHERE id = ?`, ID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`DELETE FROM bucket_tasks WHERE id = ?`, ID)

		return err
	})
}
