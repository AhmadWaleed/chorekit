package model

import (
	"database/sql"
	"strings"
)

type Task struct {
	Model
	Servers []*Server
	Runs    []*Run
	Name    string
	Env     string
	Script  string
}

func (t *Task) EnvVar() map[string]string {
	slice := strings.Split(t.Env, ";")
	vars := make(map[string]string, len(slice))
	for _, e := range vars {
		env := strings.Split(e, "=")
		vars[env[0]] = env[1]
	}

	return vars
}

type TaskServer struct {
	Model
	TaskID   uint
	ServerID uint
}

type TaskModel interface {
	Create(m *Task) error
	Update(m *Task) error
	Delete(ID uint) error
	All() ([]*Task, error)
	FindByID(ID uint) (*Task, error)
	FindMany(IDs []uint) ([]*Task, error)
}

type MysqlTaskModel struct {
	DB *sql.DB
}

func (m *MysqlTaskModel) Create(t *Task) error {
	err := trans(m.DB, func(tx *sql.Tx) error {
		result, err := tx.Exec(`INSERT INTO tasks(name, env, script) VALUES(?, ?, ?)`, t.Name, t.Env, t.Script)
		ID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for _, s := range t.Servers {
			_, err := tx.Exec(`INSERT INTO task_servers(task_id, server_id) VALUES(?, ?)`, ID, s.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (m *MysqlTaskModel) Update(t *Task) error {
	err := trans(m.DB, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			Update tasks SET 
			name=?,
			env=?,
			script=?
			WHERE id = ?`,
			t.Name,
			t.Env,
			t.Script,
			t.ID,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`DELETE FROM task_servers WHERE task_id = ?`, t.ID)
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(`INSERT INTO task_servers (task_id, server_id) VALUES(?, ?)`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, s := range t.Servers {
			_, err := stmt.Exec(t.ID, s.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (m *MysqlTaskModel) Delete(ID uint) error {
	return trans(m.DB, func(tx *sql.Tx) error {
		_, err := tx.Exec(`DELETE FROM task_servers WHERE task_id = ?`, ID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`DELETE FROM tasks WHERE id = ?`, ID)

		return err
	})
}

func (m *MysqlTaskModel) All() ([]*Task, error) {
	rows, err := m.DB.Query(`SELECT * FROM tasks`)
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for rows.Next() {
		t := new(Task)
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Env,
			&t.Script,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	rows.Close()

	for _, task := range tasks {
		var servers []*Server
		rows, err := m.DB.Query(`
			SELECT s.id, s.name 
			FROM task_servers AS ts 
			INNER JOIN servers AS s ON ts.server_id = s.id 
			WHERE ts.task_id = ?`,
			task.ID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			s := new(Server)
			if err := rows.Scan(&s.ID, &s.Name); err != nil {
				return nil, err
			}
			servers = append(servers, s)
		}
		task.Servers = servers
	}

	return tasks, rows.Err()
}

func (m *MysqlTaskModel) FindByID(ID uint) (*Task, error) {
	t := new(Task)
	row := m.DB.QueryRow(`SELECT * FROM tasks WHERE id = ?`, ID)
	if err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Env,
		&t.Script,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	rows, err := m.DB.Query(`
		SELECT s.id, s.name, s.user, s.ip, s.port, ssh_public_key, ssh_private_key
		FROM task_servers AS ts 
		INNER JOIN servers AS s ON ts.server_id = s.id 
		WHERE ts.task_id = ?`,
		t.ID,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := new(Server)
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.User,
			&s.IP,
			&s.Port,
			&s.SSHPublicKey,
			&s.SSHPrivateKey,
		); err != nil {
			return nil, err
		}
		t.Servers = append(t.Servers, s)
	}
	rows.Close()

	rows, err = m.DB.Query(`SELECT id, output FROM runs WHERE task_id = ?`, t.ID)
	if err != nil {
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		r := new(Run)
		if err := rows.Scan(&r.ID, &r.Output); err != nil {
			return nil, err
		}
		t.Runs = append(t.Runs, r)
	}

	return t, nil
}

func (m *MysqlTaskModel) FindMany(IDs []uint) ([]*Task, error) {
	var tasks []*Task
	for _, ID := range IDs {
		task, err := m.FindByID(ID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
