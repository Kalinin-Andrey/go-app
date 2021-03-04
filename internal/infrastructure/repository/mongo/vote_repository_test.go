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
	"redditclone/internal/pkg/config"

	"redditclone/internal/domain"
	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/vote"
)

type VoteRepositoryTestSuite struct {
	//	for all tests
	suite.Suite
	cfg    *config.Configuration
	logger *log.Logger
	vote   *vote.Vote
	//	only for each individual test
	ctx                context.Context
	dbMock             *dbmockmongo.DB
	voteCollectionMock *dbmockmongo.Collection
	repository         vote.Repository
}

func (s *VoteRepositoryTestSuite) SetupSuite() {
	var err error

	s.cfg = config.Get4UnitTest("CommentRepository")

	s.logger, err = log.New(s.cfg.Log)
	require.NoError(s.T(), err)

	s.vote = &vote.Vote{
		ID:        "10",
		PostID:    "1",
		UserID:    1,
		Value:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	s.dbMock = &dbmockmongo.DB{}

	s.voteCollectionMock = &dbmockmongo.Collection{}
}

func (s *VoteRepositoryTestSuite) SetupTest() {
	var ok bool
	require := require.New(s.T())
	s.ctx = context.Background()

	*s.voteCollectionMock = dbmockmongo.Collection{}
	s.dbMock.On("Collection", comment.TableName, []*options.CollectionOptions(nil)).Return(s.voteCollectionMock)

	r, err := GetRepository(s.logger, s.dbMock, comment.EntityName)
	require.NoError(err)

	s.repository, ok = r.(vote.Repository)
	require.Truef(ok, "Can not cast DB repository for entity %q to %vRepository. Repo: %v", post.EntityName, post.EntityName, r)
}

func TestVoteRepository(t *testing.T) {
	suite.Run(t, new(CommentRepositoryTestSuite))
}

func (s *VoteRepositoryTestSuite) TestGet() {
	assert := assert.New(s.T())

	result := &dbmockmongo.SingleResult{
		Entity: s.vote,
		Err:    nil,
	}

	s.voteCollectionMock.On("FindOne", s.ctx, bson.M{"id": s.vote.ID}, []*options.FindOneOptions(nil)).Return(result)

	res, err := s.repository.Get(s.ctx, s.vote.ID)
	assert.NoError(err)

	assert.Equalf(*s.vote, *res, "The two objects should be the same. Expected: %v; have got: %v", *s.vote, *res)
}

func (s *VoteRepositoryTestSuite) TestQuery() {
	var items []interface{}
	var itemsVals []vote.Vote
	assert := assert.New(s.T())

	items = append(items, s.vote)
	itemsVals = append(itemsVals, *s.vote)
	cursor := &dbmockmongo.Cursor{
		Res: items,
	}
	condition := domain.DBQueryConditions{
		Where: &vote.Vote{
			PostID: s.vote.PostID,
		},
	}
	s.voteCollectionMock.On("Find", s.ctx, bson.M{"postid": s.vote.PostID}, []*options.FindOptions(nil)).Return(cursor, error(nil))

	res, err := s.repository.Query(s.ctx, condition)
	assert.NoError(err)

	assert.Equalf(itemsVals, res, "The two objects should be the same. Expected: %v; have got: %v", itemsVals, res)
}

func (s *VoteRepositoryTestSuite) TestCreate() {
	assert := assert.New(s.T())
	newItem := &vote.Vote{}
	*newItem = *s.vote
	(*newItem).ID = ""

	s.voteCollectionMock.On("InsertOne", s.ctx, mock.Anything).Return("create test", error(nil))

	err := s.repository.Create(s.ctx, newItem)
	assert.NoError(err)
	assert.NotEmpty((*newItem).ID, "entity.ID should be is not empty")
}

func (s *VoteRepositoryTestSuite) TestUpdate() {
	assert := assert.New(s.T())

	s.voteCollectionMock.On("UpdateOne", s.ctx, bson.M{"id": s.vote.ID}, bson.M{"$set": s.vote}).Return("update test", error(nil))

	err := s.repository.Update(s.ctx, s.vote)
	assert.NoError(err)
}

func (s *VoteRepositoryTestSuite) TestDelete() {
	assert := assert.New(s.T())

	s.voteCollectionMock.On("DeleteOne", s.ctx, bson.M{"id": s.vote.ID}).Return(int64(123), error(nil))

	err := s.repository.Delete(s.ctx, s.vote.ID)
	assert.NoError(err)
}

func (s *VoteRepositoryTestSuite) TestFirst() {
	assert := assert.New(s.T())

	result := &dbmockmongo.SingleResult{
		Entity: s.vote,
		Err:    nil,
	}

	s.voteCollectionMock.On("FindOne", s.ctx, bson.M{"userid": s.vote.UserID, "postid": s.vote.PostID}, []*options.FindOneOptions(nil)).Return(result)

	res, err := s.repository.Get(s.ctx, s.vote.ID)
	assert.NoError(err)

	assert.Equalf(*s.vote, *res, "The two objects should be the same. Expected: %v; have got: %v", *s.vote, *res)
}
