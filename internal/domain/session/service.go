package session

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	// Get returns the session with the specified user ID.
	GetByUserID(ctx context.Context, userId uint) (*Session, error)
	// Create saves a new entity in the storage.
	Create(ctx context.Context, entity *Session) error
	// Update updates the entity with given ID in the storage.
	Update(ctx context.Context, entity *Session) error
	// Delete removes the entity with given ID from the storage.
	Delete(ctx context.Context, id uint) error
}

type service struct {
	//Domain     Domain
	repo       IRepository
	logger     log.ILogger
}

// NewService creates a new service.
func NewService(repo IRepository, logger log.ILogger) IService {
	s := &service{repo, logger}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s service) defaultConditions() map[string]interface{} {
	return map[string]interface{}{
	}
}

func (s service) NewEntity() *Session {
	return &Session{}
}


// Get returns the session with the specified user ID.
func (s service) GetByUserID(ctx context.Context, userId uint) (*Session, error) {
	return s.repo.GetByUserID(ctx, userId)
}

// Create saves a new entity in the storage.
func (s service) Create(ctx context.Context, entity *Session) error {
	return s.repo.Create(ctx, entity)
}

// Update updates the entity with given ID in the storage.
func (s service) Update(ctx context.Context, entity *Session) error {
	return s.repo.Create(ctx, entity)
}

// Delete removes the entity with given ID from the storage.
func (s service) Delete(ctx context.Context, id uint) (error) {
	return s.repo.Delete(ctx, id)
}

