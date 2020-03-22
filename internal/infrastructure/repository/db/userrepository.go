package db

import (
	"context"
	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"
)

// UserRepository is a repository for the user entity
type UserRepository struct {
	repository
}

var _ user.IRepository = (*UserRepository)(nil)

// New creates a new UserRepository
func NewUserRepository(repository *repository) (*UserRepository, error) {
	return &UserRepository{repository: *repository}, nil
}


// Get reads the album with the specified ID from the database.
func (r UserRepository) Get(ctx context.Context, id uint) (*user.User, error) {
	var item user.User

	r.dbWithDefaults().First(&item, id)

	return &item, nil
}

func (r UserRepository) First(ctx context.Context, entity *user.User) (*user.User, error) {
	err := r.db.DB().Where(entity).First(entity).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, apperror.ErrNotFound
		}
		r.logger.With(ctx).Error(err)
		return entity, apperror.ErrInternal
	}
	return entity, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r UserRepository) Query(ctx context.Context, offset, limit uint) ([]user.User, error) {
	var items []user.User

	r.dbWithContext(ctx, r.dbWithDefaults()).Find(&items)

	return items, nil
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r UserRepository) Create(ctx context.Context, entity *user.User) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	if err := r.db.DB().Create(entity).Error; err != nil {
		return err
	}

	return nil
}
