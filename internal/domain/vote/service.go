package vote

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
	"github.com/pkg/errors"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity(postId uint, val int) *Vote
	Get(ctx context.Context, id uint) (*Vote, error)
	Query(ctx context.Context, offset, limit uint) ([]Vote, error)
	List(ctx context.Context) ([]Vote, error)
	//Count(ctx context.Context) (uint, error)
	Vote(ctx context.Context, entity *Vote) error
	Unvote(ctx context.Context, entity *Vote) error
	Create(ctx context.Context, entity *Vote) error
	Update(ctx context.Context, entity *Vote) error
	Delete(ctx context.Context, entity *Vote) error
	First(ctx context.Context, user *Vote) (*Vote, error)
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

func (s service) NewEntity(postId uint, val int) *Vote {
	return &Vote{
		PostID:	postId,
		Value:	val,
	}
}

// Get returns the entity with the specified ID.
func (s service) Get(ctx context.Context, id uint) (*Vote, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get a vote by id: %v", id)
	}
	return entity, nil
}
/*
// Count returns the number of items.
func (s service) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit uint) ([]Vote, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of votes by ctx")
	}
	return items, nil
}

// List returns the items list.
func (s service) List(ctx context.Context) ([]Vote, error) {
	items, err := s.repo.Query(ctx, 0, MaxLIstLimit)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of votes by ctx")
	}
	return items, nil
}

func (s service) Vote(ctx context.Context, entity *Vote) (err error) {
	item := &Vote{
		PostID:		entity.PostID,
		UserID:		entity.UserID,
	}

	if item, err = s.First(ctx, item); err != nil {
		if err == apperror.ErrNotFound {
			return s.Create(ctx, entity)
		}
		return errors.Wrapf(err, "Can not find a vote by params: %v", item)
	}

	if item.Value == entity.Value {
		//	no action
		return nil
	}
	item.Value = entity.Value
	return s.Update(ctx, item)
}

func (s service) Unvote(ctx context.Context, entity *Vote) (err error) {
	item := &Vote{
		PostID:		entity.PostID,
		UserID:		entity.UserID,
	}

	if item, err = s.First(ctx, item); err != nil {
		if err == apperror.ErrNotFound {
			return err
		}
		return errors.Wrapf(err, "Can not find a vote by params: %v", item)
	}

	return s.Delete(ctx, item)
}

func (s service) Create(ctx context.Context, entity *Vote) error {
	return s.repo.Create(ctx, entity)
}

func (s service) Update(ctx context.Context, entity *Vote) error {
	return s.repo.Update(ctx, entity)
}

func (s service) Delete(ctx context.Context, entity *Vote) error {
	return s.repo.Delete(ctx, entity)
}

func (s service) First(ctx context.Context, user *Vote) (*Vote, error) {
	return s.repo.First(ctx, user)
}
