package comment

import (
	"context"

	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity() *Comment
	Get(ctx context.Context, id string) (*Comment, error)
	Query(ctx context.Context, offset, limit uint) ([]Comment, error)
	List(ctx context.Context) ([]Comment, error)
	//Count(ctx context.Context) (uint, error)
	Create(ctx context.Context, entity *Comment) error
	//Update(ctx context.Context, id string, input *Comment) (*Comment, error)
	Delete(ctx context.Context, id string) error
	//First(ctx context.Context, user *Comment) (*Comment, error)
}

type service struct {
	//Domain     Domain
	logger     log.ILogger
	repository Repository
}

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository) IService {
	s := &service{
		logger:     logger,
		repository: repo,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() selection_condition.SelectionCondition {
	return selection_condition.SelectionCondition{}
}

func (s *service) NewEntity() *Comment {
	return &Comment{}
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id string) (*Comment, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get a comment by id: %v", id)
	}
	return entity, nil
}

/*
// Count returns the number of items.
func (s *service) Count(ctx context.Context) (uint, error) {
	return s.repository.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
func (s *service) Query(ctx context.Context, offset, limit uint) ([]Comment, error) {
	items, err := s.repository.Query(ctx, selection_condition.SelectionCondition{})
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of comments by ctx")
	}
	return items, nil
}

// List returns the items list.
func (s *service) List(ctx context.Context) ([]Comment, error) {
	items, err := s.repository.Query(ctx, selection_condition.SelectionCondition{})
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of comments by ctx")
	}
	return items, nil
}

func (s *service) Create(ctx context.Context, entity *Comment) error {
	return s.repository.Create(ctx, entity)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
