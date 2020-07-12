package mock

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
)

// UserRepository is a mock for UserRepository
type CommentRepository struct {
	Response struct {
		Get struct {
			Entity *comment.Comment
			Err    error
		}
		First struct {
			Entity *comment.Comment
			Err    error
		}
		Query struct {
			List []comment.Comment
			Err  error
		}
		Create struct {
			Entity *comment.Comment
			Err    error
		}
		Update struct {
			Err error
		}
		Delete struct {
			Err error
		}
	}
}

var _ comment.IRepository = (*CommentRepository)(nil)

func (r CommentRepository) SetDefaultConditions(conditions map[string]interface{}) {}

func (r CommentRepository) Get(ctx context.Context, id uint) (*comment.Comment, error) {
	return r.Response.Get.Entity, r.Response.Get.Err
}

func (r CommentRepository) First(ctx context.Context, user *comment.Comment) (*comment.Comment, error) {
	return r.Response.First.Entity, r.Response.First.Err
}

func (r CommentRepository) Query(ctx context.Context, offset, limit uint) ([]comment.Comment, error) {
	return r.Response.Query.List, r.Response.Query.Err
}

func (r CommentRepository) Create(ctx context.Context, entity *comment.Comment) error {
	if r.Response.Create.Entity != nil {
		copier.Copy(&entity, &r.Response.Create.Entity)
	}
	return r.Response.Create.Err
}

func (r CommentRepository) Update(ctx context.Context, entity *comment.Comment) error {
	return r.Response.Update.Err
}

func (r CommentRepository) Delete(ctx context.Context, id uint) error {
	return r.Response.Delete.Err
}
