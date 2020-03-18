package db

import (
	"context"

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
func (r PostRepository) Get(ctx context.Context, id uint) (*post.Entity, error) {
	var entity post.Entity

	r.dbWithDefaults().First(&entity, id)

	return &entity, nil
}

func (r PostRepository) First(ctx context.Context, entity *post.Entity) (*post.Entity, error) {
	r.db.DB().Where(entity).First(entity)
	if entity.ID == 0 {
		return entity, apperror.ErrNotFound
	}
	return entity, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r PostRepository) Query(ctx context.Context, offset, limit uint) ([]post.Entity, error) {
	var items []post.Entity

	r.dbWithDefaults().Find(&items)

	return items, nil
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r PostRepository) Create(ctx context.Context, entity *post.Entity) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	if err := r.db.DB().Create(entity).Error; err != nil {
		return err
	}

	return nil
}
