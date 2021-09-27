package model

import (
	"database/sql"
)

type Status string

func (s Status) String() string { return string(s) }

var (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type Server struct {
	Model
	Name          string
	IP            string
	User          string
	Port          int
	SSHPublicKey  string
	SSHPrivateKey string
	Status        string
}

type ServerModel interface {
	Create(m *Server) error
	Delete(ID uint) error
	FindByID(ID uint) (*Server, error)
	FindMany(IDs []uint) ([]*Server, error)
	All() ([]*Server, error)
	UpdateStatusByID(ID uint, status Status) error
}

type MysqlServerModel struct {
	DB *sql.DB
}

func (m *MysqlServerModel) Create(s *Server) error {
	_, err := m.DB.Exec(`
		INSERT INTO servers (
		name,
		ip,
		user,
		port,
		ssh_public_key,
		ssh_private_key,
		status
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		s.Name,
		s.IP,
		s.User,
		s.Port,
		s.SSHPublicKey,
		s.SSHPrivateKey,
		s.Status,
	)
	return err
}

func (m *MysqlServerModel) Delete(ID uint) error {
	stmt, err := m.DB.Prepare(`DELETE FROM servers WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	return err
}

func (m *MysqlServerModel) FindByID(ID uint) (*Server, error) {
	s := new(Server)
	row := m.DB.QueryRow(`SELECT * FROM servers WHERE id = ?`, ID)
	if err := row.Scan(
		&s.ID,
		&s.Name,
		&s.IP,
		&s.User,
		&s.Port,
		&s.SSHPublicKey,
		&s.SSHPrivateKey,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}
	return s, row.Err()
}

func (m *MysqlServerModel) All() ([]*Server, error) {
	rows, err := m.DB.Query(`SELECT * FROM servers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []*Server
	for rows.Next() {
		s := new(Server)
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.IP,
			&s.User,
			&s.Port,
			&s.SSHPublicKey,
			&s.SSHPrivateKey,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}

	return servers, rows.Err()
}

func (m *MysqlServerModel) FindMany(IDs []uint) ([]*Server, error) {
	var servers []*Server
	for _, ID := range IDs {
		server, err := m.FindByID(ID)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}

	return servers, nil
}

func (m *MysqlServerModel) UpdateStatusByID(ID uint, status Status) error {
	_, err := m.DB.Exec(`Update servers SET status=? WHERE id = ?`, status.String(), ID)
	return err
}
