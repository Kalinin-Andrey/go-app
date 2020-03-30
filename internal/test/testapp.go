package test

import (
	"github.com/Kalinin-Andrey/redditclone/pkg/config"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"

	commonApp "github.com/Kalinin-Andrey/redditclone/internal/app"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
	"github.com/Kalinin-Andrey/redditclone/internal/test/mock"
)

// New func is a constructor for the App
func NewCommonApp(cfg config.Configuration) *commonApp.App {
	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}

	app := &commonApp.App{
		Cfg:    cfg,
		Logger: logger,
		DB:     nil,
	}

	app.Domain.User.Repository = &mock.UserRepository{}
	app.Domain.User.Service = user.NewService(app.Domain.User.Repository, app.Logger)

	app.Domain.Post.Repository = &mock.PostRepository{}
	app.Domain.Post.Service = post.NewService(app.Domain.Post.Repository, app.Logger)

	app.Domain.Vote.Repository = &mock.VoteRepository{}
	app.Domain.Vote.Service = vote.NewService(app.Domain.Vote.Repository, app.Logger)

	app.Domain.Comment.Repository = &mock.CommentRepository{}
	app.Domain.Comment.Service = comment.NewService(app.Domain.Comment.Repository, app.Logger)

	app.SessionRepository = &mock.SessionRepository{}

	return app
}

