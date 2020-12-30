package comment

import (
	"context"
	"redditclone/internal/domain"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id string) (*Comment, error)
	// Count returns the number of albums.
	//Count(ctx context.Context) (uint, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond domain.DBQueryConditions) ([]Comment, error)
	SetDefaultConditions(conditions domain.DBQueryConditions)
	// Create saves a new album in the storage.
	Create(ctx context.Context, entity *Comment) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, entity *Comment) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id string) error
	//First(ctx context.Context, user *Comment) (*Comment, error)
}
