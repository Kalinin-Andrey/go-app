package mock

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
)

// UserRepository is a mock for UserRepository
type PostRepository struct {
	Response struct {
		Get struct {
			Entity *post.Post
			Err    error
		}
		First struct {
			Entity *post.Post
			Err    error
		}
		Query struct {
			List []post.Post
			Err  error
		}
		Create struct {
			Entity *post.Post
			Err    error
		}
		Update struct {
			Entity *post.Post
			Err    error
		}
		Delete struct {
			Err error
		}
	}
}

var _ post.IRepository = (*PostRepository)(nil)

func (r PostRepository) SetDefaultConditions(conditions map[string]interface{}) {}

func (r PostRepository) Get(ctx context.Context, id uint) (*post.Post, error) {
	return r.Response.Get.Entity, r.Response.Get.Err
}

func (r PostRepository) First(ctx context.Context, user *post.Post) (*post.Post, error) {
	return r.Response.First.Entity, r.Response.First.Err
}

func (r PostRepository) Query(ctx context.Context, offset, limit uint) ([]post.Post, error) {
	return r.Response.Query.List, r.Response.Query.Err
}

func (r PostRepository) Create(ctx context.Context, entity *post.Post) error {
	if r.Response.Create.Entity != nil {
		copier.Copy(&entity, &r.Response.Create.Entity)
	}
	return r.Response.Create.Err
}

func (r PostRepository) Update(ctx context.Context, entity *post.Post) error {
	if r.Response.Create.Entity != nil {
		copier.Copy(&entity, &r.Response.Create.Entity)
	}
	return r.Response.Update.Err
}

func (r PostRepository) Delete(ctx context.Context, id uint) error {
	return r.Response.Delete.Err
}
