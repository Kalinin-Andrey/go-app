package db

import (
	"context"
	"github.com/jinzhu/gorm"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
)

// PostRepository is a repository for the user entity
type PostRepository struct {
	repository
}

var _ post.IRepository = (*PostRepository)(nil)

// New creates a new PostRepository
func NewPostRepository(repository *repository) (*PostRepository, error) {
	return &PostRepository{repository: *repository}, nil
}


// Get reads the album with the specified ID from the database.
func (r PostRepository) Get(ctx context.Context, id uint) (*post.Post, error) {
	entity := &post.Post{}

	err := r.dbWithDefaults().First(&entity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

func (r PostRepository) First(ctx context.Context, entity *post.Post) (*post.Post, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r PostRepository) Query(ctx context.Context, offset, limit uint) ([]post.Post, error) {
	var items []post.Post

	err := r.dbWithContext(ctx, r.dbWithDefaults()).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, apperror.ErrNotFound
		}
	}
	return items, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r PostRepository) Create(ctx context.Context, entity *post.Post) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

func (r PostRepository) Update(ctx context.Context, entity *post.Post) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}
	return r.db.DB().Save(entity).Error
}

// Delete deletes an entity with the specified ID from the database.
func (r PostRepository) Delete(ctx context.Context, id uint) error {
	entity, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.DB().Delete(entity).Error
}

