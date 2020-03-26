package session

import (
	"context"
)

// IRepository encapsulates the logic to access session from the data source.
type IRepository interface {
	// Get returns the session with the specified user ID.
	GetByUserID(ctx context.Context, userId uint) (*Session, error)
	// Create saves a new entity in the storage.
	Create(ctx context.Context, entity *Session) error
	// Update updates the entity with given ID in the storage.
	Update(ctx context.Context, entity *Session) error
	// Delete removes the entity with given ID from the storage.
	Delete(ctx context.Context, id uint) error
	SetDefaultConditions(conditions map[string]interface{})
	GetVar(name string) (interface{}, bool)
	SetVar(name string, val interface{}) error
}

