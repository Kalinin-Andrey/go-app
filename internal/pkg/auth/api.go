package auth

import (
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(rg *routing.RouteGroup, service Service, logger log.ILogger) {
	rg.Post("/login", login(service, logger))
	rg.Post("/register", register(service, logger))
}

func register(service Service, logger log.ILogger) routing.Handler {
	return func(c *routing.Context) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.Read(&req); err != nil {
			logger.With(c.Request.Context()).Errorf("invalid request: %v", err)
			return errorshandler.BadRequest("")
		}

		token, err := service.Register(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			if er, ok := err.(errorshandler.ErrorResponse); ok {
				logger.Errorf("Error while registering user. Status: %v; err: %q; details: %v", er.StatusCode(), er.Message, er.Details)
				return er
			}
			return err
		}
		return c.Write(struct {
			Token string `json:"token"`
		}{token})
	}
}

// login returns a handler that handles user login request.
func login(service Service, logger log.ILogger) routing.Handler {
	return func(c *routing.Context) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.Read(&req); err != nil {
			logger.With(c.Request.Context()).Errorf("invalid request: %v", err)
			return errorshandler.BadRequest("")
		}

		token, err := service.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			return err
		}
		return c.Write(struct {
			Token string `json:"token"`
		}{token})
	}
}
