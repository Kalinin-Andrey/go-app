package mongo

import (
	"context"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	mongoutil "github.com/minipkg/go-app-common/db/mongo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"redditclone/internal/pkg/apperror"

	"redditclone/internal/domain"
	"redditclone/internal/domain/vote"
)

// VoteRepository is a repository for the vote entity
type VoteRepository struct {
	repository
}

var _ vote.Repository = (*VoteRepository)(nil)

// New creates a new VoteRepository
func NewVoteRepository(repository *repository) (*VoteRepository, error) {
	return &VoteRepository{
		repository: *repository,
	}, nil
}

// Get reads the recordset with the specified ID from the database.
func (r VoteRepository) Get(ctx context.Context, id string) (*vote.Vote, error) {
	entity := &vote.Vote{}
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves records with the specified offset and limit from the database.
func (r VoteRepository) Query(ctx context.Context, cond domain.DBQueryConditions) ([]vote.Vote, error) {
	items := []vote.Vote{}
	var err error
	condition := mongoutil.QueryWhereCondition(cond.Where)

	cursor, err := r.collection.Find(ctx, condition)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return items, apperror.ErrNotFound
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "Find() error: %v", err)
	}

	for cursor.Next(ctx) {
		item := &vote.Vote{}
		cursor.Decode(item)
		items = append(items, *item)
	}
	return items, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r VoteRepository) Create(ctx context.Context, entity *vote.Vote) error {
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

func (r VoteRepository) Update(ctx context.Context, entity *vote.Vote) error {
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
func (r VoteRepository) Delete(ctx context.Context, id string) error {

	res, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return apperror.ErrNotFound
		}
		return errors.Wrapf(apperror.ErrInternal, "Can not delete entityId: %v, error: %v", id, err)
	}
	r.logger.Debugf("Delete result: %v", res)
	return nil
}

func (r VoteRepository) First(ctx context.Context, entity *vote.Vote) (*vote.Vote, error) {
	vote := &vote.Vote{}

	err := r.collection.FindOne(ctx, bson.M{"userid": entity.UserID, "postid": entity.PostID}).Decode(vote)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return vote, nil
}
