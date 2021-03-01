package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	dbmockmongo "github.com/minipkg/go-app-common/db/mongo/mock"
	"github.com/minipkg/go-app-common/log"

	"redditclone/internal/domain"
	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/post"
	"redditclone/internal/pkg/config"
)

type CommentRepositoryTestSuite struct {
	//	for all tests
	suite.Suite
	cfg     *config.Configuration
	logger  *log.Logger
	comment *comment.Comment
	//	only for each individual test
	ctx                   context.Context
	dbMock                *dbmockmongo.DB
	commentCollectionMock *dbmockmongo.Collection
	repository            comment.Repository
}

func (s *CommentRepositoryTestSuite) SetupSuite() {
	var err error

	s.cfg = config.Get4UnitTest("CommentRepository")

	s.logger, err = log.New(s.cfg.Log)
	require.NoError(s.T(), err)

	s.comment = &comment.Comment{
		ID:        "10",
		PostID:    "1",
		UserID:    1,
		Body:      "Who care about comments?",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	s.dbMock = &dbmockmongo.DB{}

	s.commentCollectionMock = &dbmockmongo.Collection{}
}

func (s *CommentRepositoryTestSuite) SetupTest() {
	var ok bool
	require := require.New(s.T())
	s.ctx = context.Background()

	*s.commentCollectionMock = dbmockmongo.Collection{}
	s.dbMock.On("Collection", comment.TableName, []*options.CollectionOptions(nil)).Return(s.commentCollectionMock)

	r, err := GetRepository(s.logger, s.dbMock, comment.EntityName)
	require.NoError(err)

	s.repository, ok = r.(comment.Repository)
	require.Truef(ok, "Can not cast DB repository for entity %q to %vRepository. Repo: %v", post.EntityName, post.EntityName, r)
}

func TestCommentRepository(t *testing.T) {
	suite.Run(t, new(CommentRepositoryTestSuite))
}

func (s *CommentRepositoryTestSuite) TestGet() {
	assert := assert.New(s.T())

	result := &dbmockmongo.SingleResult{
		Entity: s.comment,
		Err:    nil,
	}

	s.commentCollectionMock.On("FindOne", s.ctx, bson.M{"id": s.comment.ID}, []*options.FindOneOptions(nil)).Return(result)

	res, err := s.repository.Get(s.ctx, s.comment.ID)
	assert.NoError(err)

	assert.Equalf(*s.comment, *res, "The two objects should be the same. Expected: %v; have got: %v", *s.comment, *res)
}

func (s *CommentRepositoryTestSuite) TestQuery() {
	var items []interface{}
	var itemsVals []comment.Comment
	assert := assert.New(s.T())

	items = append(items, s.comment)
	itemsVals = append(itemsVals, *s.comment)
	cursor := &dbmockmongo.Cursor{
		Res: items,
	}
	condition := domain.DBQueryConditions{
		Where: map[string]interface{}{
			"PostID": s.comment.PostID,
		},
	}
	s.commentCollectionMock.On("Find", s.ctx, bson.M{"postid": s.comment.PostID}, []*options.FindOptions(nil)).Return(cursor, error(nil))

	res, err := s.repository.Query(s.ctx, condition)
	assert.NoError(err)

	assert.Equalf(itemsVals, res, "The two objects should be the same. Expected: %v; have got: %v", itemsVals, res)
}

func (s *CommentRepositoryTestSuite) TestCreate() {
	assert := assert.New(s.T())
	newItem := &comment.Comment{}
	*newItem = *s.comment
	(*newItem).ID = ""

	s.commentCollectionMock.On("InsertOne", s.ctx, mock.Anything).Return("create test", error(nil))

	err := s.repository.Create(s.ctx, newItem)
	assert.NoError(err)
	assert.NotEmpty((*newItem).ID, "entity.ID should be is not empty")
}

func (s *CommentRepositoryTestSuite) TestUpdate() {
	assert := assert.New(s.T())

	s.commentCollectionMock.On("UpdateOne", s.ctx, bson.M{"id": s.comment.ID}, bson.M{"$set": s.comment}).Return("update test", error(nil))

	err := s.repository.Update(s.ctx, s.comment)
	assert.NoError(err)
}

func (s *CommentRepositoryTestSuite) TestDelete() {
	assert := assert.New(s.T())

	s.commentCollectionMock.On("DeleteOne", s.ctx, bson.M{"id": s.comment.ID}).Return(int64(123), error(nil))

	err := s.repository.Delete(s.ctx, s.comment.ID)
	assert.NoError(err)
}
