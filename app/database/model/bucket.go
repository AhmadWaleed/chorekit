package model

import (
	"database/sql"
)

type Bucket struct {
	Model
	Name     string
	Parallel bool
	Tasks    []*BucketTask
	Runs     []*BucketRuns
}

type BucketTask struct {
	Model
	Task     *Task
	TaskID   uint
	BucketID uint
}

type BucketRuns struct {
	BucketID  uint
	TaskRunID uint
	Run       *Run
}

type BucketModel interface {
	Create(b *Bucket) error
	CreateRun(b *Bucket, r *Run) error
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

func (m *MysqlBucketModel) CreateRun(b *Bucket, r *Run) error {
	return trans(m.DB, func(tx *sql.Tx) error {
		result, err := tx.Exec(`INSERT INTO task_runs (task_id, output) VALUES (?, ?)`, r.TaskID, r.Output)
		if err != nil {
			return err
		}
		ID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		_, err = tx.Exec(`INSERT INTO bucket_runs (bucket_id, task_run_id) VALUES (?, ?)`, b.ID, ID)
		return err
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
	defer rows.Close()

	for _, b := range buckets {
		if err := m.Load(b, "tasks"); err != nil {
			return buckets, err
		}
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

	if err := m.Load(b, "tasks"); err != nil {
		return b, err
	}

	if err := m.Load(b, "runs"); err != nil {
		return b, err
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

func (m *MysqlBucketModel) Load(b *Bucket, rel string) error {
	switch rel {
	case "tasks":
		rows, err := m.DB.Query(`
			SELECT 
			bt.id,
			bt.task_id,
			t.id,
			t.name,
			t.env,
			t.script
			FROM bucket_tasks AS bt 
			INNER JOIN tasks AS t ON bt.task_id = t.id 
			WHERE bt.bucket_id = ?`,
			b.ID,
		)
		if err != nil {
			return err
		}

		var tasks []*BucketTask
		for rows.Next() {
			t := new(Task)
			bt := new(BucketTask)
			if err := rows.Scan(
				&bt.ID,
				&bt.TaskID,
				&t.ID,
				&t.Name,
				&t.Env,
				&t.Script,
			); err != nil {
				return err
			}
			bt.Task = t
			tasks = append(tasks, bt)
		}
		b.Tasks = tasks

		return rows.Close()
	case "runs":
		rows, err := m.DB.Query(`
			SELECT 
			br.bucket_id,
			br.task_run_id,
			tr.id,
			tr.task_id,
			tr.output,
			tr.created_at
			FROM bucket_runs AS br 
			INNER JOIN task_runs AS tr ON br.task_run_id = tr.id 
			WHERE br.bucket_id = ?`,
			b.ID,
		)
		if err != nil {
			return err
		}

		var runs []*BucketRuns
		for rows.Next() {
			r := new(Run)
			br := new(BucketRuns)
			if err := rows.Scan(
				&br.BucketID,
				&br.TaskRunID,
				&r.ID,
				&r.TaskID,
				&r.Output,
				&r.CreatedAt,
			); err != nil {
				return err
			}
			br.Run = r
			runs = append(runs, br)
		}
		b.Runs = runs

		return rows.Close()
	}

	return nil
}
