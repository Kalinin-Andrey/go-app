package mock

import (
	"context"

	"github.com/Kalinin-Andrey/redditclone/internal/pkg/session"
)

// UserRepository is a mock for UserRepository
type SessionRepository struct {
	Response struct {
		NewEntity		struct {
			Entity	*session.Session
			Err		error
		}
		GetByUserID	struct {
			Entity	*session.Session
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
		SetVar	struct {
			Err		error
		}
		GetVar	struct {
			Val		interface{}
			Ok		bool
		}
	}
}

var _ session.IRepository = (*SessionRepository)(nil)

func (r SessionRepository) SetDefaultConditions(conditions map[string]interface{}) {}

func (r SessionRepository) NewEntity(ctx context.Context, userId uint) (*session.Session, error) {
	return r.Response.NewEntity.Entity, r.Response.NewEntity.Err
}

func (r SessionRepository) GetByUserID(ctx context.Context, userId uint) (*session.Session, error) {
	return r.Response.GetByUserID.Entity, r.Response.GetByUserID.Err
}

func (r SessionRepository) Create(ctx context.Context, entity *session.Session) error {
	return r.Response.Create.Err
}

func (r SessionRepository) Update(ctx context.Context, entity *session.Session) error {
	return r.Response.Update.Err
}

func (r SessionRepository) Delete(ctx context.Context, id uint) error {
	return r.Response.Delete.Err
}

func (r SessionRepository) GetVar(session *session.Session, name string) (interface{}, bool) {
	return r.Response.GetVar.Val, r.Response.GetVar.Ok
}

func (r SessionRepository) SetVar(session *session.Session, name string, val interface{}) error {
	return r.Response.SetVar.Err
}
