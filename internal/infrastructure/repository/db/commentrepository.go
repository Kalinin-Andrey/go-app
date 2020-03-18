package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
)

// CommentRepository is a repository for the user entity
type CommentRepository struct {
	repository
}

var _ comment.IRepository = (*CommentRepository)(nil)

// New creates a new CommentRepository
func NewCommentRepository(repository *repository) (*CommentRepository, error) {
	return &CommentRepository{repository: *repository}, nil
}


// Get reads the album with the specified ID from the database.
func (r CommentRepository) Get(ctx context.Context, id uint) (*comment.Entity, error) {
	var entity comment.Entity

	r.dbWithDefaults().First(&entity, id)

	return &entity, nil
}

func (r CommentRepository) First(ctx context.Context, entity *comment.Entity) (*comment.Entity, error) {
	r.db.DB().Where(entity).First(entity)
	if entity.ID == 0 {
		return entity, apperror.ErrNotFound
	}
	return entity, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r CommentRepository) Query(ctx context.Context, offset, limit uint) ([]comment.Entity, error) {
	var items []comment.Entity

	r.dbWithDefaults().Find(&items)

	return items, nil
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r CommentRepository) Create(ctx context.Context, entity *comment.Entity) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	if err := r.db.DB().Create(entity).Error; err != nil {
		return err
	}

	return nil
}
