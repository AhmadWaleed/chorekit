package model

import (
	"database/sql"
	"strings"
	"time"
)

type Run struct {
	ID        uint
	Task      Task
	TaskID    uint
	Output    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Run) DisplayOutput() []string {
	return strings.Split(r.Output, "\n")
}

type RunModel interface {
	Create(m *Run) error
	FindByID(ID uint) (*Run, error)
	FindByTaskID(ID uint) ([]*Run, error)
}

type MysqlRunModel struct {
	DB *sql.DB
}

func (m *MysqlRunModel) Create(r *Run) error {
	_, err := m.DB.Exec(`INSERT INTO runs (task_id, output) VALUES (?, ?)`, r.TaskID, r.Output)
	return err
}

func (m *MysqlRunModel) FindByID(ID uint) (*Run, error) {
	r := new(Run)
	row := m.DB.QueryRow(`SELECT * FROM runs WHERE id = ?`, ID)
	if err := row.Scan(
		&r.ID,
		&r.TaskID,
		&r.Output,
		&r.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}

		return nil, err
	}
	return r, row.Err()
}

func (m *MysqlRunModel) FindByTaskID(ID uint) ([]*Run, error) {
	rows, err := m.DB.Query(`SELECT * FROM runs WHERE task_id = ?`, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []*Run
	for rows.Next() {
		r := new(Run)
		if err := rows.Scan(
			&r.ID,
			&r.TaskID,
			&r.Output,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		runs = append(runs, r)
	}

	return runs, rows.Err()
}
