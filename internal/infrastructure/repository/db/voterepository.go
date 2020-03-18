package db

import (
	"context"

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
func (r VoteRepository) Get(ctx context.Context, id uint) (*vote.Entity, error) {
	var entity vote.Entity

	r.dbWithDefaults().First(&entity, id)

	return &entity, nil
}

func (r VoteRepository) First(ctx context.Context, entity *vote.Entity) (*vote.Entity, error) {
	r.db.DB().Where(entity).First(entity)
	if entity.ID == 0 {
		return entity, apperror.ErrNotFound
	}
	return entity, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r VoteRepository) Query(ctx context.Context, offset, limit uint) ([]vote.Entity, error) {
	var items []vote.Entity

	r.dbWithDefaults().Find(&items)

	return items, nil
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r VoteRepository) Create(ctx context.Context, entity *vote.Entity) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	if err := r.db.DB().Create(entity).Error; err != nil {
		return err
	}

	return nil
}
