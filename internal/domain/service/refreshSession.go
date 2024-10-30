package service

import "auth-api/internal/domain/entity"

type RefreshSessionStorage interface {
	GetByIDUser(id int) (sessions []entity.RefreshSession, err error)
	Create(session entity.RefreshSession) error
}

type refreshSessionService struct {
	refreshSessionStorage *RefreshSessionStorage
}

func NewRefreshSessionService(refreshSessionStorage *RefreshSessionStorage) *refreshSessionService {
	return &refreshSessionService{refreshSessionStorage: refreshSessionStorage}
}
