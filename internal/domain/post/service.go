package post

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity() *Post
	Get(ctx context.Context, id uint) (*Post, error)
	First(ctx context.Context, user *Post) (*Post, error)
	Query(ctx context.Context, offset, limit uint) ([]Post, error)
	List(ctx context.Context) ([]Post, error)
	//Count(ctx context.Context) (uint, error)
	Create(ctx context.Context, entity *Post) error
	ViewsIncr(ctx context.Context, entity *Post) error
	//Update(ctx context.Context, entity *Post) error
	Delete(ctx context.Context, id uint) (error)
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

func (s service) NewEntity() *Post {
	return &Post{}
}

// Get returns the entity with the specified ID.
func (s service) Get(ctx context.Context, id uint) (*Post, error) {
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
func (s service) Query(ctx context.Context, offset, limit uint) ([]Post, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return items, nil
}


// List returns the items list.
func (s service) List(ctx context.Context) ([]Post, error) {
	items, err := s.repo.Query(ctx, 0, MaxLIstLimit)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s service) Create(ctx context.Context, entity *Post) error {
	return s.repo.Create(ctx, entity)
}

func (s service) ViewsIncr(ctx context.Context, entity *Post) error {
	entity.Views++
	return s.repo.Update(ctx, entity)
}

func (s service) First(ctx context.Context, user *Post) (*Post, error) {
	return s.repo.First(ctx, user)
}

func (s service) Delete(ctx context.Context, id uint) (error) {
	return s.repo.Delete(ctx, id)
}
