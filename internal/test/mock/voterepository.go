package mock

import (
	"context"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
)

// UserRepository is a mock for UserRepository
type VoteRepository struct {
	Response struct {
		Get		struct {
			Entity	*vote.Vote
			Err		error
		}
		First	struct {
			Entity	*vote.Vote
			Err		error
		}
		Query	struct {
			List	[]vote.Vote
			Err		error
		}
		Create	struct {
			Err		error
		}
		Update	struct {
			Err		error
		}
		Delete	struct {
			Err		error
		}
	}
}

var _ vote.IRepository = (*VoteRepository)(nil)

func (r VoteRepository) SetDefaultConditions(conditions map[string]interface{}) {}

func (r VoteRepository) Get(ctx context.Context, id uint) (*vote.Vote, error) {
	return r.Response.Get.Entity, r.Response.Get.Err
}

func (r VoteRepository) First(ctx context.Context, user *vote.Vote) (*vote.Vote, error) {
	return r.Response.First.Entity, r.Response.First.Err
}

func (r VoteRepository) Query(ctx context.Context, offset, limit uint) ([]vote.Vote, error) {
	return r.Response.Query.List, r.Response.Query.Err
}

func (r VoteRepository) Create(ctx context.Context, entity *vote.Vote) error {
	return r.Response.Create.Err
}

func (r VoteRepository) Update(ctx context.Context, entity *vote.Vote) error {
	return r.Response.Update.Err
}

func (r VoteRepository) Delete(ctx context.Context, entity *vote.Vote) error {
	return r.Response.Delete.Err
}
