package post

import (
	"context"
)

// IRepository encapsulates the logic to access albums from the data source.
type IRepository interface {
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*Entity, error)
	// Count returns the number of albums.
	//Count(ctx context.Context) (uint, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, offset, limit uint) ([]Entity, error)
	SetDefaultConditions(conditions map[string]interface{})
	// Create saves a new album in the storage.
	Create(ctx context.Context, entity *Entity) error
	// Update updates the album with given ID in the storage.
	//Update(ctx context.Context, entity *Entity) error
	// Delete removes the album with given ID from the storage.
	//Delete(ctx context.Context, id uint) error
	First(ctx context.Context, user *Entity) (*Entity, error)
}

