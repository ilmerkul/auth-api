package mysql

import (
	"auth-api/internal/domain/entity"
	"context"
	"database/sql"
	"time"
)

const (
	contextTimeGetSession    = 2
	contextTimeCreateSession = 1
)

type refreshSessionStorage struct {
	db      *sql.DB
	context context.Context
}

func NewRefreshSessionStorage(db *sql.DB) *refreshSessionStorage {
	return &refreshSessionStorage{db: db, context: context.Background()}
}

func (s *refreshSessionStorage) GetByIDUser(id int) (sessions []entity.RefreshSession, err error) {
	q := `SELECT * FROM refreshSession WHERE userId=?`

	context, close := context.WithTimeout(s.context, contextTimeGetSession*time.Second)
	defer close()

	if err = s.db.PingContext(context); err != nil {
		return
	}

	stmt, err := s.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	session := entity.RefreshSession{}
	rows, err := stmt.QueryContext(context, id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&session); err != nil {
			return
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *refreshSessionStorage) Create(session entity.RefreshSession) error {

	q := `INSERT INTO refreshSession (userId, refreshToken, ua, fingerprint, expiresIn, createdAt) VALUES (?, ?, ?, ?, ?, ?)`

	context, close := context.WithTimeout(s.context, contextTimeCreateSession*time.Second)
	defer close()

	stmt, err := s.db.PrepareContext(context, q)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(context, session.UserId, session.RefreshToken, session.UA, session.Fingerprint, session.ExpiresIn, session.CreatedAt)

	return err
}
