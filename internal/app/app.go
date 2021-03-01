package app

import (
	"context"
	"github.com/minipkg/go-app-common/log"
	"github.com/pkg/errors"
	golog "log"
	"redditclone/internal/pkg/apperror"
	"redditclone/internal/pkg/config"

	"github.com/minipkg/go-app-common/db/mongo"
	"github.com/minipkg/go-app-common/db/pg"
	"github.com/minipkg/go-app-common/db/redis"
	"github.com/minipkg/go-app-common/db/redis/cache"
	"redditclone/internal/pkg/auth"
	"redditclone/internal/pkg/jwt"

	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/user"
	"redditclone/internal/domain/vote"
	mongorep "redditclone/internal/infrastructure/repository/mongo"
	pgrep "redditclone/internal/infrastructure/repository/pg"
	redisrep "redditclone/internal/infrastructure/repository/redis"
)

// App struct is the common part of all applications
type App struct {
	Cfg     config.Configuration
	Logger  log.ILogger
	DB      pg.IDB
	MongoDB mongo.IDB
	Redis   redis.IDB
	Domain  Domain
	Auth    Auth
	Cache   cache.Service
}

type Auth struct {
	SessionRepository auth.SessionRepository
	TokenRepository   auth.TokenRepository
	Service           auth.Service
}

// Domain is a Domain Layer Entry Point
type Domain struct {
	User    DomainUser
	Post    DomainPost
	Vote    DomainVote
	Comment DomainComment
}

type DomainUser struct {
	Repository user.Repository
	Service    user.IService
}

type DomainPost struct {
	Repository post.Repository
	Service    post.IService
}

type DomainVote struct {
	Repository vote.Repository
	Service    vote.IService
}

type DomainComment struct {
	Repository comment.Repository
	Service    comment.IService
}

// New func is a constructor for the App
func New(cfg config.Configuration) *App {
	logger, err := log.New(cfg.Log)
	if err != nil {
		golog.Fatal(err)
	}

	pgDB, err := pg.New(cfg.DB.Pg, logger)
	if err != nil {
		golog.Fatal(err)
	}

	mDB, err := mongo.New(cfg.DB.Mongo)
	if err != nil {
		golog.Fatal(err)
	}

	rDB, err := redis.New(cfg.DB.Redis)
	if err != nil {
		golog.Fatal(err)
	}

	app := &App{
		Cfg:     cfg,
		Logger:  logger,
		DB:      pgDB,
		MongoDB: mDB,
		Redis:   rDB,
	}

	err = app.Init()
	if err != nil {
		golog.Fatal(err)
	}

	return app
}

func (app *App) Init() (err error) {
	if err := app.SetupRepositories(); err != nil {
		return err
	}
	app.SetupServices()
	return nil
}

func (app *App) SetupRepositories() (err error) {
	var ok bool

	app.Domain.User.Repository, ok = app.getPgRepo(user.EntityName).(user.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", user.EntityName, user.EntityName, app.getPgRepo(user.EntityName))
	}

	app.Domain.Post.Repository, ok = app.getMongoRepo(post.EntityName).(post.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", post.EntityName, post.EntityName, app.getMongoRepo(post.EntityName))
	}

	app.Domain.Vote.Repository, ok = app.getMongoRepo(vote.EntityName).(vote.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", vote.EntityName, vote.EntityName, app.getMongoRepo(post.EntityName))
	}

	app.Domain.Comment.Repository, ok = app.getMongoRepo(comment.EntityName).(comment.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", comment.EntityName, comment.EntityName, app.getMongoRepo(post.EntityName))
	}

	if app.Auth.SessionRepository, err = redisrep.NewSessionRepository(app.Redis, app.Cfg.SessionLifeTime, app.Domain.User.Repository); err != nil {
		return errors.Errorf("Can not get new SessionRepository err: %v", err)
	}
	app.Auth.TokenRepository = jwt.NewRepository()

	app.Cache = cache.NewService(app.Redis, app.Cfg.CacheLifeTime)

	return nil
}

func (app *App) SetupServices() {
	app.Domain.User.Service = user.NewService(app.Logger, app.Domain.User.Repository)
	app.Domain.Post.Service = post.NewService(app.Logger, app.Domain.Post.Repository, app.Domain.Comment.Repository, app.Domain.Vote.Repository)
	app.Domain.Vote.Service = vote.NewService(app.Logger, app.Domain.Vote.Repository)
	app.Domain.Comment.Service = comment.NewService(app.Logger, app.Domain.Comment.Repository)
	app.Auth.Service = auth.NewService(app.Cfg.JWTSigningKey, app.Cfg.JWTExpiration, app.Domain.User.Service, app.Logger, app.Auth.SessionRepository, app.Auth.TokenRepository)
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}

func (app *App) getPgRepo(entityName string) (repo pgrep.IRepository) {
	var err error

	if repo, err = pgrep.GetRepository(app.Logger, app.DB, entityName); err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entityName, err)
	}
	return repo
}

func (app *App) getMongoRepo(entityName string) (repo mongorep.IRepository) {
	var err error

	if repo, err = mongorep.GetRepository(app.Logger, app.MongoDB, entityName); err != nil {
		golog.Fatalf("Can not get mongodb repository for entity %q, error happened: %v", entityName, err)
	}
	return repo
}

func (app *App) Stop() error {
	errRedis := app.Redis.Close()
	errPg := app.DB.DB().Close()
	errMongo := app.MongoDB.Close(context.Background())

	switch {
	case errPg != nil:
		return errors.Wrapf(apperror.ErrInternal, "pg error: %v", errPg)
	case errMongo != nil:
		return errors.Wrapf(apperror.ErrInternal, "mongo error: %v", errMongo)
	case errRedis != nil:
		return errors.Wrapf(apperror.ErrInternal, "redis error: %v", errRedis)
	}

	return nil
}
