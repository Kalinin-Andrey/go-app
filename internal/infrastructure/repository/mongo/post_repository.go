package mongo

import (
	"context"
	"github.com/pkg/errors"
	"strings"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"redditclone/internal/pkg/apperror"

	"redditclone/internal/domain"
	"redditclone/internal/domain/post"
)

// PostRepository is a repository for the post entity
type PostRepository struct {
	repository
	commentRepository *CommentRepository
	voteRepository    *VoteRepository
}

var _ post.Repository = (*PostRepository)(nil)

// New creates a new PostRepository
func NewPostRepository(repository *repository, commentRepository *CommentRepository, voteRepository *VoteRepository) (*PostRepository, error) {
	return &PostRepository{
		repository:        *repository,
		commentRepository: commentRepository,
		voteRepository:    voteRepository,
	}, nil
}

// Get reads the recordset with the specified ID from the database.
func (r PostRepository) Get(ctx context.Context, id string) (*post.Post, error) {
	entity := &post.Post{}
	res := r.collection.FindOne(ctx, bson.M{"id": id})
	err := res.Decode(entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotFound
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "FindOne() error: %v", err)
	}

	err = r.populate(ctx, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (r PostRepository) populate(ctx context.Context, item *post.Post) (err error) {
	cond := domain.DBQueryConditions{
		Where: map[string]interface{}{"PostID": item.ID},
	}

	comments, err := r.commentRepository.Query(ctx, cond)
	if err != nil {
		return err
	}

	votes, err := r.voteRepository.Query(ctx, cond)
	if err != nil {
		return err
	}

	(*item).Comments = comments
	(*item).Votes = votes
	return nil
}

// Query retrieves records with the specified offset and limit from the database.
func (r PostRepository) Query(ctx context.Context, cond domain.DBQueryConditions) ([]post.Post, error) {
	var err error
	items := []post.Post{}
	condition	:= bson.M{}

	for k, v := range cond.Where {
		condition[strings.ToLower(k)] = v
	}
	// bson.M{"postid": cond.Where["PostID"].(string)}

	cursor, err := r.collection.Find(ctx, condition)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return items, apperror.ErrNotFound
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "Find() error: %v", err)
	}

	for cursor.Next(ctx) {
		item := &post.Post{}
		cursor.Decode(item)

		err = r.populate(ctx, item)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}
	return items, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r PostRepository) Create(ctx context.Context, entity *post.Post) error {
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

func (r PostRepository) Update(ctx context.Context, entity *post.Post) error {
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
func (r PostRepository) Delete(ctx context.Context, id string) error {
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
