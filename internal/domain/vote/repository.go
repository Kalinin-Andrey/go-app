package vote

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id string) (*Vote, error)
	// Count returns the number of albums.
	//Count(ctx context.Context) (uint, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond selection_condition.SelectionCondition) ([]Vote, error)
	// Create saves a new album in the storage.
	Create(ctx context.Context, entity *Vote) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, entity *Vote) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id string) error
	First(ctx context.Context, entity *Vote) (*Vote, error)
}
