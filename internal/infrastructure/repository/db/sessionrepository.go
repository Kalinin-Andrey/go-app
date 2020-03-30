package db

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/db"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/session"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"
)

// PostRepository is a repository for the user entity
type SessionRepository struct {
	repository
	UserRepo user.IRepository
}

var _ session.IRepository = (*SessionRepository)(nil)

// New creates a new PostRepository
func NewSessionRepository(dbase db.IDB, logger log.ILogger, userRepo user.IRepository) (*SessionRepository, error) {
	r := &SessionRepository{
		repository: repository{
			db:     dbase,
			logger: logger,
		},
		UserRepo: userRepo,
	}

	return r, nil
}

func (r SessionRepository) NewEntity(ctx context.Context, userId uint) (*session.Session, error) {
	user, err := r.UserRepo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &session.Session{
		UserID:		userId,
		User:		*user,
		Data:		make(map[string]interface{}, 1),
		Json:		"{}",
	}, nil
}

func (r SessionRepository) GetVar(session *session.Session, name string) (interface{}, bool) {

	if session.Data == nil {
		session.Data = make(map[string]interface{}, 1)
	}

	val, ok := session.Data[name]
	return val, ok
}

func (r *SessionRepository) SetVar(session *session.Session, name string, val interface{}) error {

	if session.Data == nil {
		session.Data = make(map[string]interface{}, 1)
	}

	session.Data[name] = val
	return r.SaveSession(session)
}

func (r *SessionRepository) SaveSession(session *session.Session) error {
	return r.Update(session.Ctx, session)
}


// Get reads the album with the specified ID from the database.
func (r SessionRepository) Get(ctx context.Context, id uint) (*session.Session, error) {
	var entity session.Session

	err := r.dbWithDefaults().First(&entity, id).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &entity, apperror.ErrNotFound
		}
		r.logger.With(ctx).Error(err)
		return &entity, apperror.ErrInternal
	}
	entity.SetDataByJson()

	return &entity, err
}

// Get returns the Session with the specified user ID.
func (r SessionRepository) GetByUserID(ctx context.Context, userId uint) (*session.Session, error) {
	var entity session.Session

	err := r.dbWithDefaults().Where(&session.Session{
		UserID: userId,
	}).First(&entity).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &entity, apperror.ErrNotFound
		}
		r.logger.With(ctx).Error(err)
		return &entity, apperror.ErrInternal
	}
	entity.SetDataByJson()

	return &entity, err
}

// Create saves a new entity in the storage.
func (r SessionRepository) Create(ctx context.Context, entity *session.Session) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	entity.SetJsonByData()
	return r.db.DB().Create(entity).Error
}

// Update updates the entity with given ID in the storage.
func (r SessionRepository) Update(ctx context.Context, entity *session.Session) error {
	entity.SetJsonByData()
	return r.db.DB().Save(entity).Error
}

// Delete removes the entity with given ID from the storage.
func (r SessionRepository) Delete(ctx context.Context, id uint) (error) {
	entity, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	return r.db.DB().Delete(entity).Error
}

