package comment

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
	"github.com/pkg/errors"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity() *Comment
	Get(ctx context.Context, id uint) (*Comment, error)
	Query(ctx context.Context, offset, limit uint) ([]Comment, error)
	List(ctx context.Context) ([]Comment, error)
	//Count(ctx context.Context) (uint, error)
	Create(ctx context.Context, entity *Comment) error
	//Update(ctx context.Context, id string, input *Comment) (*Comment, error)
	Delete(ctx context.Context, id uint) (error)
	First(ctx context.Context, user *Comment) (*Comment, error)
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

func (s service) NewEntity() *Comment {
	return &Comment{}
}

// Get returns the entity with the specified ID.
func (s service) Get(ctx context.Context, id uint) (*Comment, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get a comment by id: %v", id)
	}
	return entity, nil
}
/*
// Count returns the number of items.
func (s service) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit uint) ([]Comment, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of comments by ctx")
	}
	return items, nil
}

// List returns the items list.
func (s service) List(ctx context.Context) ([]Comment, error) {
	items, err := s.repo.Query(ctx, 0, MaxLIstLimit)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of comments by ctx")
	}
	return items, nil
}

func (s service) Create(ctx context.Context, entity *Comment) error {
	return s.repo.Create(ctx, entity)
}

func (s service) First(ctx context.Context, user *Comment) (*Comment, error) {
	return s.repo.First(ctx, user)
}

func (s service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
