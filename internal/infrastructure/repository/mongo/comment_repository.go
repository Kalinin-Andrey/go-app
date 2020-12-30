package mongo

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"redditclone/internal/pkg/apperror"

	"redditclone/internal/domain"
	"redditclone/internal/domain/comment"
)

// CommentRepository is a repository for the comment entity
type CommentRepository struct {
	repository
}

var _ comment.Repository = (*CommentRepository)(nil)

// New creates a new CommentRepository
func NewCommentRepository(repository *repository) (*CommentRepository, error) {
	return &CommentRepository{
		repository: *repository,
	}, nil
}

// Get reads the recordset with the specified ID from the database.
func (r CommentRepository) Get(ctx context.Context, id string) (*comment.Comment, error) {
	entity := &comment.Comment{}
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves records with the specified offset and limit from the database.
func (r CommentRepository) Query(ctx context.Context, cond domain.DBQueryConditions) ([]comment.Comment, error) {
	items := []comment.Comment{}
	var err error

	cursor, err := r.collection.Find(ctx, bson.M{"postid": cond.Where["PostID"].(string)})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return items, apperror.ErrNotFound
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "Find() error: %v", err)
	}

	for cursor.Next(ctx) {
		item := &comment.Comment{}
		err = cursor.Decode(item)
		if err != nil {
			return nil, errors.Wrapf(apperror.ErrInternal, "Decode() error: %v", err)
		}
		items = append(items, *item)
	}
	return items, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r CommentRepository) Create(ctx context.Context, entity *comment.Comment) error {
	if entity.ID != "" {
		return errors.Wrap(apperror.ErrBadRequest, "entity is not new")
	}

	entity.ID = uuid.New().String()

	id, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Can not create a recordset for an object %v, error: %v", entity, err)
	}
	r.logger.Debugf("Create records InsertedID: %v", id)
	return nil
}

func (r CommentRepository) Update(ctx context.Context, entity *comment.Comment) error {
	if entity.ID == "" {
		return errors.Wrap(apperror.ErrBadRequest, "entity is new")
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"id": entity.ID}, bson.M{"$set": entity})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return apperror.ErrNotFound
		}
		return errors.Wrapf(apperror.ErrInternal, "Can not update entity: %v, error: %v", entity, err)
	}
	r.logger.Debugf("Update result: %v", res)

	return nil
}

// Delete deletes an entity with the specified ID from the database.
func (r CommentRepository) Delete(ctx context.Context, id string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return apperror.ErrNotFound
		}
		return errors.Wrapf(apperror.ErrInternal, "Can not delete entity id: %v, error: %v", id, err)
	}
	r.logger.Debugf("Delete result: %v", res)
	return nil
}
