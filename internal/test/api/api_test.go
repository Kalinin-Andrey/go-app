package api

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"redditclone/internal/pkg/session"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"redditclone/internal/pkg/config"
	"redditclone/internal/pkg/jwt"
	repositoryMock "redditclone/internal/pkg/mock/repository"

	"github.com/minipkg/go-app-common/log"

	commonapp "redditclone/internal/app"
	apiapp "redditclone/internal/app/restapi"
	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/user"
	"redditclone/internal/domain/vote"
)

type ApiTestSuite struct {
	//	for all tests
	suite.Suite
	cfg      *config.Configuration
	logger   *log.Logger
	api      *apiapp.App
	server   *httptest.Server
	client   *http.Client
	token    string
	entities entities
	//	only for each individual test
	ctx             context.Context
	repositoryMocks repositoryMocks
}

type entities struct {
	user    *user.User
	post    *post.Post
	comment *comment.Comment
	vote    *vote.Vote
}

type repositoryMocks struct {
	user    *repositoryMock.UserRepository
	session *repositoryMock.SessionRepository
	post    *repositoryMock.PostRepository
	comment *repositoryMock.CommentRepository
	vote    *repositoryMock.VoteRepository
}

func (s *ApiTestSuite) SetupSuite() {
	var err error

	s.cfg = config.Get4UnitTest("api")

	s.logger, err = log.New(s.cfg.Log)
	require.NoError(s.T(), err)

	s.client = &http.Client{}

	s.setupEntities()
	s.initMocks()
}

func (s *ApiTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.setupMocks()
	s.api = apiapp.New(s.newCommonApp(), *s.cfg)

	s.server = httptest.NewServer(s.api.Server.Handler)
}

func (s *ApiTestSuite) AfterTest(_, _ string) {
	s.server.Close()
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (s *ApiTestSuite) newCommonApp() *commonapp.App {
	app := &commonapp.App{
		Cfg:    *s.cfg,
		Logger: s.logger,
	}
	app.Domain.User.Repository = s.repositoryMocks.user
	app.Domain.Post.Repository = s.repositoryMocks.post
	app.Domain.Comment.Repository = s.repositoryMocks.comment
	app.Domain.Vote.Repository = s.repositoryMocks.vote
	app.Auth.SessionRepository = s.repositoryMocks.session
	app.Auth.TokenRepository = jwt.NewRepository()

	app.SetupServices()
	return app
}

func (s *ApiTestSuite) setupEntities() {
	var passhash, _ = hex.DecodeString("3a73acfdb534ddded4c0109383ee3e5a66314113d1ff691aaf4b3ee073c8fc2edd06d48f0555ec3783f4c479994e3eee3433734c29b05f08be0e9739b956b88d8fe872bd0a0942214e94fd4001e757fa3b66a2b9925de2e800c55ef49baa4c03")
	s.entities = entities{
		user: &user.User{
			ID:        1,
			Name:      "demo1",
			Passhash:  string(passhash),
			CreatedAt: time.Now().Local(),
			UpdatedAt: time.Now().Local(),
		},
		vote: &vote.Vote{
			ID:        "12",
			PostID:    "1",
			UserID:    1,
			Value:     1,
			CreatedAt: time.Now().Local(),
			UpdatedAt: time.Now().Local(),
			DeletedAt: nil,
		},
	}
	s.entities.comment = &comment.Comment{
		ID:        "11",
		PostID:    "1",
		UserID:    1,
		Body:      "Who care about comments?",
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		DeletedAt: nil,
		User:      *s.entities.user,
	}
	s.entities.post = &post.Post{
		ID:        "1",
		Score:     2,
		Views:     3,
		Title:     "What does a good programmer mean?",
		Type:      post.TypeText,
		Category:  post.CategoryProgramming,
		Text:      "Who can consider himself a good programmer?",
		UserID:    1,
		User:      *s.entities.user,
		Votes:     []vote.Vote{*s.entities.vote},
		Comments:  []comment.Comment{*s.entities.comment},
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		DeletedAt: nil,
	}
}

func (s *ApiTestSuite) initMocks() {
	s.repositoryMocks = repositoryMocks{
		user:    &repositoryMock.UserRepository{},
		session: &repositoryMock.SessionRepository{},
		post:    &repositoryMock.PostRepository{},
		comment: &repositoryMock.CommentRepository{},
		vote:    &repositoryMock.VoteRepository{},
	}
}

func (s *ApiTestSuite) setupMocks() {
	*s.repositoryMocks.user = repositoryMock.UserRepository{}
	*s.repositoryMocks.session = repositoryMock.SessionRepository{}
	*s.repositoryMocks.post = repositoryMock.PostRepository{}
	*s.repositoryMocks.comment = repositoryMock.CommentRepository{}
	*s.repositoryMocks.vote = repositoryMock.VoteRepository{}
}

func (s *ApiTestSuite) setupSession() {
	newSession := &session.Session{
		UserID: s.entities.user.ID,
		User:   *s.entities.user,
		Data: session.Data{
			UserID:              s.entities.user.ID,
			UserName:            s.entities.user.Name,
			ExpirationTokenTime: time.Now().Local().Add(time.Hour),
		},
	}
	s.repositoryMocks.session.On("Get", mock.Anything, s.entities.user.ID).Return(newSession, error(nil))
}
