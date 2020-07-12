package app

import (
	golog "log"

	"github.com/Kalinin-Andrey/redditclone/pkg/config"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"

	"github.com/Kalinin-Andrey/redditclone/internal/pkg/db"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/session"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
	dbrep "github.com/Kalinin-Andrey/redditclone/internal/infrastructure/repository/db"
)

type IApp interface {
	// Run is func to run the App
	Run() error
}

// App struct is the common part of all applications
type App struct {
	Cfg               config.Configuration
	Logger            log.ILogger
	DB                db.IDB
	Domain            Domain
	SessionRepository session.IRepository
}

// Domain is a Domain Layer Entry Point
type Domain struct {
	User struct {
		Repository user.IRepository
		Service    user.IService
	}
	Post struct {
		Repository post.IRepository
		Service    post.IService
	}
	Vote struct {
		Repository vote.IRepository
		Service    vote.IService
	}
	Comment struct {
		Repository comment.IRepository
		Service    comment.IService
	}
}

const (
	EntityNameUser    = "user"
	EntityNamePost    = "post"
	EntityNameVote    = "vote"
	EntityNameComment = "comment"
)

// New func is a constructor for the App
func New(cfg config.Configuration) *App {
	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}

	db, err := db.New(cfg.DB, logger)
	if err != nil {
		panic(err)
	}

	app := &App{
		Cfg:    cfg,
		Logger: logger,
		DB:     db,
	}
	var ok bool

	app.Domain.User.Repository, ok = app.getDBRepo(EntityNameUser).(user.IRepository)
	if !ok {
		golog.Fatalf("Can not cast DB repository for entity %q to %v.IRepository. Repo: %v", EntityNameUser, EntityNameUser, app.getDBRepo(EntityNameUser))
	}
	app.Domain.User.Service = user.NewService(app.Domain.User.Repository, app.Logger)

	app.Domain.Post.Repository, ok = app.getDBRepo(EntityNamePost).(post.IRepository)
	if !ok {
		golog.Fatalf("Can not cast DB repository for entity %q to %v.IRepository. Repo: %v", EntityNamePost, EntityNamePost, app.getDBRepo(EntityNamePost))
	}
	app.Domain.Post.Service = post.NewService(app.Domain.Post.Repository, app.Logger)

	app.Domain.Vote.Repository, ok = app.getDBRepo(EntityNameVote).(vote.IRepository)
	if !ok {
		golog.Fatalf("Can not cast DB repository for entity %q to %v.IRepository. Repo: %v", EntityNameVote, EntityNameVote, app.getDBRepo(EntityNameVote))
	}
	app.Domain.Vote.Service = vote.NewService(app.Domain.Vote.Repository, app.Logger)

	app.Domain.Comment.Repository, ok = app.getDBRepo(EntityNameComment).(comment.IRepository)
	if !ok {
		golog.Fatalf("Can not cast DB repository for entity %q to %v.IRepository. Repo: %v", EntityNameComment, EntityNameComment, app.getDBRepo(EntityNameComment))
	}
	app.Domain.Comment.Service = comment.NewService(app.Domain.Comment.Repository, app.Logger)

	if app.SessionRepository, err = dbrep.NewSessionRepository(app.DB, app.Logger, app.Domain.User.Repository); err != nil {
		golog.Fatalf("Can not get new SessionRepository err: %v", err)
	}

	return app
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}

func (app *App) getDBRepo(entityName string) (repo dbrep.IRepository) {
	var err error

	if repo, err = dbrep.GetRepository(app.DB, app.Logger, entityName); err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entityName, err)
	}
	return repo
}
