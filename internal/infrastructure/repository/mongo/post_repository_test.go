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

	dbmockmongo "github.com/minipkg/db/mongo/mock"
	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"

	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/user"
	"redditclone/internal/domain/vote"
	"redditclone/internal/pkg/config"
)

type PostRepositoryTestSuite struct {
	//	for all tests
	suite.Suite
	cfg     *config.Configuration
	logger  *log.Logger
	post    *post.Post
	comment *comment.Comment
	vote    *vote.Vote
	//	only for each individual test
	ctx                   context.Context
	dbMock                *dbmockmongo.DB
	postCollectionMock    *dbmockmongo.Collection
	commentCollectionMock *dbmockmongo.Collection
	voteCollectionMock    *dbmockmongo.Collection
	repository            post.Repository
}

func (s *PostRepositoryTestSuite) SetupSuite() {
	var err error

	s.cfg = config.Get4UnitTest("PostRepository")

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
	s.vote = &vote.Vote{
		ID:        "10",
		PostID:    "1",
		UserID:    1,
		Value:     1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
	s.post = &post.Post{
		ID:       "1",
		Score:    2,
		Views:    3,
		Title:    "What does a good programmer mean?",
		Type:     post.TypeText,
		Category: post.CategoryProgramming,
		Text:     "Who can consider himself a good programmer?",
		UserID:   1,
		User: user.User{
			ID:        1,
			Name:      "demo1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
		Votes:     []vote.Vote{*s.vote},
		Comments:  []comment.Comment{*s.comment},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	s.dbMock = &dbmockmongo.DB{}

	s.postCollectionMock = &dbmockmongo.Collection{}
	s.commentCollectionMock = &dbmockmongo.Collection{}
	s.voteCollectionMock = &dbmockmongo.Collection{}
}

func (s *PostRepositoryTestSuite) SetupTest() {
	var ok bool
	require := require.New(s.T())
	s.ctx = context.Background()

	*s.postCollectionMock = dbmockmongo.Collection{}
	*s.commentCollectionMock = dbmockmongo.Collection{}
	*s.voteCollectionMock = dbmockmongo.Collection{}
	s.dbMock.On("Collection", post.TableName, []*options.CollectionOptions(nil)).Return(s.postCollectionMock)
	s.dbMock.On("Collection", comment.TableName, []*options.CollectionOptions(nil)).Return(s.commentCollectionMock)
	s.dbMock.On("Collection", vote.TableName, []*options.CollectionOptions(nil)).Return(s.voteCollectionMock)

	r, err := GetRepository(s.logger, s.dbMock, post.EntityName)
	require.NoError(err)

	s.repository, ok = r.(post.Repository)
	require.Truef(ok, "Can not cast DB repository for entity %q to %vRepository. Repo: %v", post.EntityName, post.EntityName, r)
}

func TestPostRepository(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}

func (s *PostRepositoryTestSuite) populatePost() {
	var comments []interface{}
	var votes []interface{}

	comments = append(comments, s.comment)
	votes = append(votes, s.vote)
	commentCursor := &dbmockmongo.Cursor{
		Res: comments,
	}
	voteCursor := &dbmockmongo.Cursor{
		Res: votes,
	}

	s.commentCollectionMock.On("Find", s.ctx, bson.M{"postid": s.post.ID}, []*options.FindOptions(nil)).Return(commentCursor, error(nil))
	s.voteCollectionMock.On("Find", s.ctx, bson.M{"postid": s.post.ID}, []*options.FindOptions(nil)).Return(voteCursor, error(nil))
}

func (s *PostRepositoryTestSuite) TestGet() {
	assert := assert.New(s.T())

	result := &dbmockmongo.SingleResult{
		Entity: s.post,
		Err:    nil,
	}

	s.populatePost()
	s.postCollectionMock.On("FindOne", s.ctx, bson.M{"id": s.post.ID}, []*options.FindOneOptions(nil)).Return(result)

	res, err := s.repository.Get(s.ctx, s.post.ID)
	assert.NoError(err)

	assert.Equalf(*s.post, *res, "The two objects should be the same. Expected: %v; have got: %v", *s.post, *res)
}

func (s *PostRepositoryTestSuite) TestQuery() {
	var posts []interface{}
	var postVals []post.Post
	assert := assert.New(s.T())

	posts = append(posts, s.post)
	postVals = append(postVals, *s.post)
	cursor := &dbmockmongo.Cursor{
		Res: posts,
	}
	condition := selection_condition.SelectionCondition{
		Where: &post.Post{
			UserID: s.post.UserID,
		},
	}
	s.populatePost()
	s.postCollectionMock.On("Find", s.ctx, bson.M{"userid": s.post.UserID}, []*options.FindOptions(nil)).Return(cursor, error(nil))

	res, err := s.repository.Query(s.ctx, condition)
	assert.NoError(err)

	assert.Equalf(postVals, res, "The two objects should be the same. Expected: %v; have got: %v", postVals, res)
}

func (s *PostRepositoryTestSuite) TestCreate() {
	assert := assert.New(s.T())
	newPost := &post.Post{}
	*newPost = *s.post
	(*newPost).ID = ""

	s.postCollectionMock.On("InsertOne", s.ctx, mock.Anything).Return("create test", error(nil))

	err := s.repository.Create(s.ctx, newPost)
	assert.NoError(err)
	assert.NotEmpty((*newPost).ID, "entity.ID should be is not empty")
}

func (s *PostRepositoryTestSuite) TestUpdate() {
	assert := assert.New(s.T())

	s.postCollectionMock.On("UpdateOne", s.ctx, bson.M{"id": s.post.ID}, bson.M{"$set": s.post}).Return("update test", error(nil))

	err := s.repository.Update(s.ctx, s.post)
	assert.NoError(err)
}

func (s *PostRepositoryTestSuite) TestDelete() {
	assert := assert.New(s.T())

	s.postCollectionMock.On("DeleteOne", s.ctx, bson.M{"id": s.post.ID}).Return(int64(123), error(nil))

	err := s.repository.Delete(s.ctx, s.post.ID)
	assert.NoError(err)
}
