package api

import (
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/auth"
	"log"
	"net/http"
	"time"

	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/go-ozzo/ozzo-routing/v2/file"
	"github.com/go-ozzo/ozzo-routing/v2/slash"

	"github.com/Kalinin-Andrey/redditclone/pkg/accesslog"
	"github.com/Kalinin-Andrey/redditclone/pkg/config"
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"

	commonApp "github.com/Kalinin-Andrey/redditclone/internal/app"
	"github.com/Kalinin-Andrey/redditclone/internal/app/api/controller"
)

// Version of API
const Version = "1.0.0"

// App is the application for API
type App struct {
	*commonApp.App
	Server		*http.Server
}

// New func is a constructor for the ApiApp
func New(cfg config.Configuration) *App {
	app := &App{
		App: commonApp.New(cfg),
		Server:	nil,
	}

	// build HTTP server
	server := &http.Server{
		Addr:		cfg.Server.HTTPListen,
		Handler:	app.buildHandler(),
	}
	app.Server = server

	return app
}

func (app *App) buildHandler() *routing.Router {
	router := routing.New()

	router.Use(
		accesslog.Handler(app.Logger),
		slash.Remover(http.StatusMovedPermanently),
		errorshandler.Handler(app.Logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	// serve index file
	router.Get("/", file.Content("website/index.html"))
	router.Get("/", file.Content("website/favicon.ico"))
	router.Get("/", file.Content("website/manifest.json"))
	// serve files under the "ui" subdirectory
	router.Get("/css/*", file.Server(file.PathMap{
		"/css/": "/website/css/",
	}))
	router.Get("/js/*", file.Server(file.PathMap{
		"/js/": "/website/js/",
	}))
	rg := router.Group("/api")

	authHandler := auth.Handler(app.Cfg.JWTSigningKey, app.DB, app.Logger, app.Domain.User.Repository)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(app.Cfg.JWTSigningKey, app.Cfg.JWTExpiration, app.Domain.User.Service, app.Logger),
		app.Logger,
	)

	app.RegisterHandlers(rg, authHandler)

	return 	router
}

// Run is func to run the ApiApp
func (app *App) Run() error {
	go func() {
		defer func() {
			if err := app.DB.DB().Close(); err != nil {
				app.Logger.Error(err)
			}

			err := app.Logger.Sync()
			if err != nil {
				log.Println(err.Error())
			}
		}()
		// start the HTTP server with graceful shutdown
		routing.GracefulShutdown(app.Server, 10*time.Second, app.Logger.Infof)
	}()
	app.Logger.Infof("server %v is running at %v", Version, app.Server.Addr)
	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func (app *App) RegisterHandlers(rg *routing.RouteGroup, authHandler routing.Handler) {

	//controller.RegisterUserHandlers(rg, app.Domain.User.Service, app.Logger, authHandler)
	controller.RegisterPostHandlers(rg, app.Domain.Post.Service, app.Domain.User.Service, app.Logger, authHandler)


}
