package post

import (
	"context"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
	"go.mongodb.org/mongo-driver/mongo"

	"redditclone/internal/domain"
	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/vote"
	"redditclone/internal/pkg/apperror"
)

const MaxLIstLimit = 1000

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity() *Post
	NewVoteEntity(userId uint, postId string, val int) *vote.Vote
	Get(ctx context.Context, id string) (*Post, error)
	//First(ctx context.Context, user *Post) (*Post, error)
	Query(ctx context.Context, query domain.DBQueryConditions) ([]Post, error)
	List(ctx context.Context) ([]Post, error)
	//Count(ctx context.Context) (uint, error)
	Create(ctx context.Context, entity *Post) error
	ViewsIncr(ctx context.Context, entity *Post) error
	//Update(ctx context.Context, entity *Post) error
	Delete(ctx context.Context, id string) error
	Vote(ctx context.Context, entity *vote.Vote) error
	Unvote(ctx context.Context, entity *vote.Vote) error
}

type service struct {
	//Domain     Domain
	logger            log.ILogger
	repository        Repository
	commentRepository comment.Repository
	voteReporitory    vote.Repository
}

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository, commentRepo comment.Repository, voteRepo vote.Repository) IService {
	s := &service{
		logger:            logger,
		repository:        repo,
		commentRepository: commentRepo,
		voteReporitory:    voteRepo,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() domain.DBQueryConditions {
	return domain.DBQueryConditions{}
}

func (s *service) NewEntity() *Post {
	return &Post{}
}

func (s *service) NewVoteEntity(userId uint, postId string, val int) *vote.Vote {
	return &vote.Vote{
		UserID: userId,
		PostID: postId,
		Value:  val,
	}
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id string) (*Post, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

/*
// Count returns the number of items.
func (s *service) Count(ctx context.Context) (uint, error) {
	return s.repository.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
func (s *service) Query(ctx context.Context, query domain.DBQueryConditions) ([]Post, error) {
	items, err := s.repository.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of posts by query: %v", query)
	}
	return items, nil
}

// List returns the items list.
func (s *service) List(ctx context.Context) ([]Post, error) {
	items, err := s.repository.Query(ctx, domain.DBQueryConditions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of posts by ctx")
	}
	return items, nil
}

func (s *service) Create(ctx context.Context, entity *Post) error {
	return s.repository.Create(ctx, entity)
}

func (s *service) ViewsIncr(ctx context.Context, entity *Post) error {
	entity.Views++
	return s.repository.Update(ctx, entity)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) Vote(ctx context.Context, entity *vote.Vote) (err error) {
	item := &vote.Vote{
		PostID: entity.PostID,
		UserID: entity.UserID,
	}

	if item, err = s.voteReporitory.First(ctx, item); err != nil {
		if err == apperror.ErrNotFound {
			err = s.voteReporitory.Create(ctx, entity)
			return s.PostChangeScore(ctx, entity.PostID, entity.Value)
		}
		return errors.Wrapf(err, "Can not find a vote by params: %v", item)
	}

	if item.Value == entity.Value {
		//	no action
		return nil
	}
	item.Value = entity.Value
	if err = s.voteReporitory.Update(ctx, item); err != nil {
		return err
	}
	return s.PostChangeScore(ctx, entity.PostID, 2*entity.Value)
}

func (s *service) Unvote(ctx context.Context, entity *vote.Vote) (err error) {
	item := &vote.Vote{
		PostID: entity.PostID,
		UserID: entity.UserID,
	}

	if item, err = s.voteReporitory.First(ctx, item); err != nil {
		if err == apperror.ErrNotFound {
			return err
		}
		return errors.Wrapf(err, "Can not find a vote by params: %v", item)
	}

	if err = s.voteReporitory.Delete(ctx, item.ID); err != nil {
		return err
	}
	return s.PostChangeScore(ctx, entity.PostID, -1*entity.Value)
}

func (s *service) PostChangeScore(ctx context.Context, id string, diff int) error {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.Wrapf(apperror.ErrNotFound, "Post id: %q not found", id)
		}
		return errors.Wrapf(apperror.ErrInternal, "Post id: %q not found", id)
	}
	entity.Score += diff

	err = s.repository.Update(ctx, entity)
	if err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Can not update post: %v, error: %v", entity, err)
	}
	return nil
}
