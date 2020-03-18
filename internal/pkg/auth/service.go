package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) (string, error)
	NewUser(username, password string) (*user.Entity, error)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetName() string
}

type UserService interface {
}

type service struct {
	signingKey      string
	tokenExpiration int
	userService		user.IService
	logger          log.ILogger
}

const (
	saltSize   = 64
	iterations = 1e4
)

// NewService creates a new authentication service.
func NewService(signingKey string, tokenExpiration int, userService user.IService, logger log.ILogger) Service {
	return service{signingKey, tokenExpiration, userService, logger}
}

func (s service) NewUser(username, password string) (*user.Entity, error) {
	user := s.userService.NewEntity()
	user.Name = username

	salt, err := generateRandomBytes(saltSize)
	if err != nil {
		return user, errors.Wrapf(err, "could not get salt: %v", err)
	}
	user.Passhash = string(hashPassword([]byte(password), salt))
	return user, nil
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.authenticate(ctx, username, password)

	if err != nil {
		return "", err
	}
	return s.generateJWT(user)
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an *user.Entity is returned. Otherwise, error is returned.
func (s service) authenticate(ctx context.Context, username, password string) (*user.Entity, error) {
	logger := s.logger.With(ctx, "user", username)

	user := s.userService.NewEntity()
	user.Name = username

	user, err := s.userService.First(ctx, user)
	if err != nil {
		return user, errorshandler.BadRequest("User not found")
	}

	if comparePassword([]byte(user.Passhash), []byte(password)) {
		logger.Infof("authentication successful")
		return user, nil
	}

	logger.Infof("authentication failed")
	return user, errorshandler.Unauthorized("")
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(user *user.Entity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}

func (s service) Register(ctx context.Context, username, password string) (string, error) {
	user, err := s.NewUser(username, password)
	if err != nil {
		return "", errorshandler.InternalServerError(err.Error())
	}

	if err := s.userService.Create(ctx, user); err != nil {
		return "", errorshandler.BadRequest(err.Error())
	}

	return s.generateJWT(user)
}


// Source: https://play.golang.org/p/tAZtO7L6pm
// hash provided clear text password and compare it to provided hash
func comparePassword(hash, pw []byte) bool {
	return bytes.Equal(hash, hashPassword(pw, hash[:saltSize]))
}

// hash the password with the provided salt using the pbkdf2 algorithm
// return byte slice containing salt (first 64 bytes) and hash (last 32 bytes) => total of 96 bytes
func hashPassword(pw, salt []byte) []byte {
	ret := make([]byte, len(salt))
	copy(ret, salt)
	return append(ret, pbkdf2.Key(pw, salt, iterations, sha256.Size, sha256.New)...)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

