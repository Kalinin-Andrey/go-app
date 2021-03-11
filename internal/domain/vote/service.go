package vote

import (
	"context"
	"redditclone/internal/domain"

	"github.com/minipkg/go-app-common/log"
	"github.com/pkg/errors"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity(userId uint, postId string, val int) *Vote
	Get(ctx context.Context, id string) (*Vote, error)
	Query(ctx context.Context, offset, limit uint) ([]Vote, error)
	List(ctx context.Context) ([]Vote, error)
	//Count(ctx context.Context) (uint, error)
	Create(ctx context.Context, entity *Vote) error
	Update(ctx context.Context, entity *Vote) error
	Delete(ctx context.Context, id string) error
	First(ctx context.Context, user *Vote) (*Vote, error)
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
func (s *service) defaultConditions() domain.DBQueryConditions {
	return domain.DBQueryConditions{}
}

func (s *service) NewEntity(userId uint, postId string, val int) *Vote {
	return &Vote{
		UserID: userId,
		PostID: postId,
		Value:  val,
	}
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id string) (*Vote, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get a vote by id: %v", id)
	}
	return entity, nil
}

/*
// Count returns the number of items.
func (s *service) Count(ctx context.Context) (uint, error) {
	return s.repository.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
func (s *service) Query(ctx context.Context, offset, limit uint) ([]Vote, error) {
	items, err := s.repository.Query(ctx, domain.DBQueryConditions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of votes by ctx")
	}
	return items, nil
}

// List returns the items list.
func (s *service) List(ctx context.Context) ([]Vote, error) {
	items, err := s.repository.Query(ctx, domain.DBQueryConditions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of votes by ctx")
	}
	return items, nil
}

func (s *service) First(ctx context.Context, entity *Vote) (*Vote, error) {
	return s.repository.First(ctx, entity)
}

func (s *service) Create(ctx context.Context, entity *Vote) error {
	return s.repository.Create(ctx, entity)
}

func (s *service) Update(ctx context.Context, entity *Vote) error {
	return s.repository.Update(ctx, entity)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
