package user

import (
	"context"
	"database/sql"
)

type Store interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetById(ctx context.Context, id int) (User, error)
	UpdateName(ctx context.Context, user User) (User, error)
	//Delete(ctx context.Context, id string) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(ctx context.Context, user *User) (*User, error) {
	row := s.db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id", user.Name, user.Mail, user.Password)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return user, nil
}

func (s *PostgresStore) GetById(ctx context.Context, id int) (User, error) {
	var user User
	row := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Mail, &user.Password, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *PostgresStore) UpdateName(ctx context.Context, user User) (User, error) {
	res, err := s.db.Exec("UPDATE users SET name = $1 WHERE id = $2", user.Name, user.ID)
	if err != nil {
		return User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return User{}, err
	}
	user.ID = id
	return user, nil
}
