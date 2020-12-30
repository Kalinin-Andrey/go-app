package auth

import (
	"context"
	"redditclone/internal/domain"
	"redditclone/internal/pkg/session"
)

// SessionRepository encapsulates the logic to access session from the data source.
type SessionRepository interface {
	SetDefaultConditions(defaultConditions domain.DBQueryConditions)
	NewEntity(ctx context.Context, userId uint) (*session.Session, error)
	// GetByUserID returns the session with the specified user ID.
	Get(ctx context.Context, userId uint) (*session.Session, error)
	// Create saves a new entity in the storage.
	Create(ctx context.Context, entity *session.Session) error
	// Update updates the entity with given ID in the storage.
	Update(ctx context.Context, entity *session.Session) error
	Save(session *session.Session) error
	// Delete removes the entity with given ID from the storage.
	Delete(ctx context.Context, entity *session.Session) error
	GetData(session *session.Session) session.Data
	SetData(session *session.Session, data session.Data) error
}
