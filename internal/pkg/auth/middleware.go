package auth

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/db"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/session"
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
)

type contextKey int

const (
	userSessionKey contextKey = iota
)

// Handler returns a JWT-based authentication middleware.
//func Handler(verificationKey string, dbase db.IDB, logger log.ILogger) routing.Handler {
func Handler(verificationKey string, dbase db.IDB, logger log.ILogger, sessRepo session.IRepository) routing.Handler {
	return auth.JWT(verificationKey, auth.JWTOptions{TokenHandler: func(c *routing.Context, token *jwt.Token) error {
		ctx := c.Request.Context()
		userID := uint(token.Claims.(jwt.MapClaims)["id"].(float64))

		session, err := sessRepo.GetByUserID(ctx, userID)
		if err != nil {

			if err == apperror.ErrNotFound {
				if session, err = sessRepo.NewEntity(ctx, userID); err != nil {
					logger.With(ctx).Error(err)
					return errorshandler.InternalServerError("")
				}

				if err := sessRepo.Create(ctx, session); err != nil {
					logger.With(ctx).Error(err)
					return errorshandler.InternalServerError("")
				}
			}
		}
		session.Ctx = ctx

		ctx = context.WithValue(
			ctx,
			userSessionKey,
			session,
		)
		c.Request = c.Request.WithContext(ctx)
		return nil
	}})
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
/*func handleToken(c *routing.Context, token *jwt.Token) error {
	ctx := WithUser(
		c.Request.Context(),
		uint(token.Claims.(jwt.MapClaims)["id"].(float64)),
		token.Claims.(jwt.MapClaims)["name"].(string),
	)
	c.Request = c.Request.WithContext(ctx)
	return nil
}

// WithUser returns a context that contains the user identity from the given JWT.
func WithUser(ctx context.Context, id uint, name string) context.Context {
	return context.WithValue(ctx, userSessionKey, user.User{ID: id, Name: name})
}*/

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentSession(ctx context.Context) *session.Session {
	if sess, ok := ctx.Value(userSessionKey).(*session.Session); ok {
		return sess
	}
	return nil
}

// MockAuthHandler creates a mock authentication middleware for testing purpose.
// If the request contains an Authorization header whose value is "TEST", then
// it considers the user is authenticated as "Tester" whose ID is "100".
// It fails the authentication otherwise.
/*func MockAuthHandler(c *routing.Context) error {
	if c.Request.Header.Get("Authorization") != "TEST" {
		return errorshandler.Unauthorized("")
	}
	ctx := WithUser(c.Request.Context(), 100, "Tester")
	c.Request = c.Request.WithContext(ctx)
	return nil
}

// MockAuthHeader returns an HTTP header that can pass the authentication check by MockAuthHandler.
func MockAuthHeader() http.Header {
	header := http.Header{}
	header.Add("Authorization", "TEST")
	return header
}*/
