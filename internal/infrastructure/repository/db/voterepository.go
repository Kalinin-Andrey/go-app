package db

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
)

// VoteRepository is a repository for the user entity
type VoteRepository struct {
	repository
}

var _ vote.IRepository = (*VoteRepository)(nil)

// New creates a new VoteRepository
func NewVoteRepository(repository *repository) (*VoteRepository, error) {
	return &VoteRepository{repository: *repository}, nil
}


// Get reads the album with the specified ID from the database.
func (r VoteRepository) Get(ctx context.Context, id uint) (*vote.Vote, error) {
	entity := &vote.Vote{}

	err := r.dbWithDefaults().First(entity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

func (r VoteRepository) First(ctx context.Context, entity *vote.Vote) (*vote.Vote, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r VoteRepository) Query(ctx context.Context, offset, limit uint) ([]vote.Vote, error) {
	var items []vote.Vote

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
func (r VoteRepository) Create(ctx context.Context, entity *vote.Vote) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	tx := r.db.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(entity).Error; err != nil {
		tx.Rollback()
		return err
	}
	post := &post.Post{ID: entity.PostID}
	if err := tx.Model(post).UpdateColumn("score", gorm.Expr("score + ?", entity.Value)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}


func (r VoteRepository) Update(ctx context.Context, entity *vote.Vote) error {
	var d int = 2 * entity.Value

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}

	tx := r.db.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//if err := tx.Save(entity).Error; err != nil {
	if err := tx.Model(entity).UpdateColumn("value", entity.Value).Error; err != nil {
		tx.Rollback()
		return err
	}
	post := &post.Post{ID: entity.PostID}
	if err := tx.Model(post).UpdateColumn("score", gorm.Expr("score + ?", d)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}


func (r VoteRepository) Delete(ctx context.Context, entity *vote.Vote) error {
	var d int = -1 * entity.Value

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}

	tx := r.db.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Delete(entity).Error; err != nil {
		tx.Rollback()
		return err
	}
	post := &post.Post{ID: entity.PostID}
	if err := tx.Model(post).UpdateColumn("score", gorm.Expr("score + ?", d)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
