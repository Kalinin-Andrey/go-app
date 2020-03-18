package comment

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity() *Entity
	Get(ctx context.Context, id uint) (*Entity, error)
	Query(ctx context.Context, offset, limit uint) ([]Entity, error)
	List(ctx context.Context) ([]Entity, error)
	//Count(ctx context.Context) (uint, error)
	Create(ctx context.Context, entity *Entity) error
	//Update(ctx context.Context, id string, input *Entity) (*Entity, error)
	//Delete(ctx context.Context, id string) (error)
	First(ctx context.Context, user *Entity) (*Entity, error)
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

func (s service) NewEntity() *Entity {
	return &Entity{}
}

// Get returns the entity with the specified ID.
func (s service) Get(ctx context.Context, id uint) (*Entity, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
/*
// Count returns the number of items.
func (s service) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit uint) ([]Entity, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// List returns the items list.
func (s service) List(ctx context.Context) ([]Entity, error) {
	items, err := s.repo.Query(ctx, 0, MaxLIstLimit)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s service) Create(ctx context.Context, entity *Entity) error {
	return s.repo.Create(ctx, entity)
}

func (s service) First(ctx context.Context, user *Entity) (*Entity, error) {
	return s.repo.First(ctx, user)
}
