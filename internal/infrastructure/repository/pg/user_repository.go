package pg

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"redditclone/internal/pkg/apperror"

	"redditclone/internal/domain/user"
)

// UserRepository is a repository for the user entity
type UserRepository struct {
	repository
}

var _ user.Repository = (*UserRepository)(nil)

// New creates a new UserRepository
func NewUserRepository(repository *repository) (*UserRepository, error) {
	return &UserRepository{repository: *repository}, nil
}

// Get reads the album with the specified ID from the database.
func (r UserRepository) Get(ctx context.Context, id uint) (*user.User, error) {
	entity := &user.User{}

	err := r.dbWithDefaults().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

func (r UserRepository) First(ctx context.Context, entity *user.User) (*user.User, error) {
	err := r.dbWithDefaults().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r UserRepository) Query(ctx context.Context, offset, limit uint) ([]user.User, error) {
	items := []user.User{}

	err := r.dbWithDefaults().Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, apperror.ErrNotFound
		}
	}
	return items, err
}

// Create saves a new user record in the database.
func (r UserRepository) Create(ctx context.Context, entity *user.User) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}
