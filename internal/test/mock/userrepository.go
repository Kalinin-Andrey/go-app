package mock

import (
	"context"
	"github.com/jinzhu/copier"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
)

// UserRepository is a mock for UserRepository
type UserRepository struct {
	Response struct {
		Get		struct {
			Entity	*user.User
			Err		error
		}
		First	struct {
			Entity	*user.User
			Err		error
		}
		Query	struct {
			List	[]user.User
			Err		error
		}
		Create	struct {
			Entity	*user.User
			Err		error
		}
	}
}

var _ user.IRepository = (*UserRepository)(nil)

func (r UserRepository) SetDefaultConditions(conditions map[string]interface{}) {}

func (r UserRepository) Get(ctx context.Context, id uint) (*user.User, error) {
	return r.Response.Get.Entity, r.Response.Get.Err
}

func (r UserRepository) First(ctx context.Context, user *user.User) (*user.User, error) {
	return r.Response.First.Entity, r.Response.First.Err
}

func (r UserRepository) Query(ctx context.Context, offset, limit uint) ([]user.User, error) {
	return r.Response.Query.List, r.Response.Query.Err
}

func (r UserRepository) Create(ctx context.Context, entity *user.User) error {
	if r.Response.Create.Entity != nil {
		copier.Copy(&entity, &r.Response.Create.Entity)
	}
	return r.Response.Create.Err
}
