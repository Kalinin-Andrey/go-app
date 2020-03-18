package db

import (
	"context"

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
func (r UserRepository) Get(ctx context.Context, id uint) (*user.Entity, error) {
	var item user.Entity

	r.dbWithDefaults().First(&item, id)

	return &item, nil
}

func (r UserRepository) First(ctx context.Context, user *user.Entity) (*user.Entity, error) {
	r.db.DB().Where(user).First(user)
	if user.ID == 0 {
		return user, apperror.ErrNotFound
	}
	return user, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r UserRepository) Query(ctx context.Context, offset, limit uint) ([]user.Entity, error) {
	var items []user.Entity

	r.dbWithDefaults().Find(&items)

	return items, nil
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r UserRepository) Create(ctx context.Context, entity *user.Entity) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	if err := r.db.DB().Create(entity).Error; err != nil {
		return err
	}

	return nil
}
