package db

import (
	"context"
	"github.com/jinzhu/gorm"

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
func (r CommentRepository) Get(ctx context.Context, id uint) (*comment.Comment, error) {
	entity := &comment.Comment{}

	err := r.dbWithDefaults().First(&entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

func (r CommentRepository) First(ctx context.Context, entity *comment.Comment) (*comment.Comment, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r CommentRepository) Query(ctx context.Context, offset, limit uint) ([]comment.Comment, error) {
	var items []comment.Comment

	err := r.dbWithContext(ctx, r.dbWithDefaults()).Find(&items).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return items, apperror.ErrNotFound
		}
	}
	return items, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r CommentRepository) Create(ctx context.Context, entity *comment.Comment) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Delete deletes an entity with the specified ID from the database.
func (r CommentRepository) Delete(ctx context.Context, id uint) error {
	entity, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.DB().Delete(entity).Error
}
