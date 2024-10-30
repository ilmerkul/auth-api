package mysql

import (
	"auth-api/internal/domain/entity"
	"context"
	"database/sql"
	"time"
)

const (
	contextTimeGetOneUser = 1
	contextTimeCreateUser = 1
)

type userStorage struct {
	db      *sql.DB
	context context.Context
}

func NewUserStorage(db *sql.DB) *userStorage {
	return &userStorage{db: db, context: context.Background()}
}

func (s *userStorage) GetByID(id int) (user entity.User, err error) {
	q := `SELECT * FROM user WHERE id=?`

	context, close := context.WithTimeout(s.context, contextTimeGetOneUser*time.Second)
	defer close()

	if err = s.db.PingContext(context); err != nil {
		return
	}

	stmt, err := s.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	user = entity.User{}
	if err = stmt.QueryRowContext(context, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return
	}

	return user, nil
}

func (s *userStorage) Create(user entity.User) (id int64, err error) {

	q := `INSERT INTO user (email, password, createdAt) VALUES (?, ?, ?)`

	context, close := context.WithTimeout(s.context, contextTimeCreateUser*time.Second)
	defer close()

	stmt, err := s.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	row, err := stmt.ExecContext(context, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return
	}

	return row.LastInsertId()
}
