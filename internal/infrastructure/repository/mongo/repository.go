package mongo

import (
	"github.com/pkg/errors"

	mongodb "github.com/minipkg/db/mongo"
	"github.com/minipkg/log"

	"redditclone/internal/domain"
	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/vote"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	logger     log.ILogger
	Conditions domain.DBQueryConditions
	db         mongodb.IDB
	collection mongodb.ICollection
}

const DefaultLimit = 100

// GetRepository return a repository
func GetRepository(logger log.ILogger, db mongodb.IDB, entity string) (repo IRepository, err error) {
	r := &repository{
		logger: logger,
		db:     db,
	}

	switch entity {
	case post.EntityName:
		r.collection = r.db.Collection(post.TableName)

		commentRepository, err := NewCommentRepository(&repository{
			logger:     logger,
			db:         db,
			collection: r.db.Collection(comment.TableName),
		})
		if err != nil {
			return nil, err
		}

		voteRepository, err := NewVoteRepository(&repository{
			logger:     logger,
			db:         db,
			collection: r.db.Collection(vote.TableName),
		})
		if err != nil {
			return nil, err
		}

		repo, err = NewPostRepository(r, commentRepository, voteRepository)
	case vote.EntityName:
		r.collection = r.db.Collection(vote.TableName)
		repo, err = NewVoteRepository(r)
	case comment.EntityName:
		r.collection = r.db.Collection(comment.TableName)
		repo, err = NewCommentRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entity)
	}
	return repo, err
}

func (r *repository) SetDefaultConditions(defaultConditions domain.DBQueryConditions) {
	r.Conditions = defaultConditions

	if r.Conditions.Limit == 0 {
		r.Conditions.Limit = DefaultLimit
	}
}
