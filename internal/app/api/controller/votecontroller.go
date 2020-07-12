package controller

import (
	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/auth"
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
	"github.com/go-ozzo/ozzo-routing/v2"
	"net/http"
)

type voteController struct {
	Controller
	Service     vote.IService
	PostService post.IService
	Logger      log.ILogger
}

//	GET /api/post/{POST_ID}/upvote - рейтинг постп вверх
//	GET /api/post/{POST_ID}/downvote - рейтинг поста вниз
//	GET /api/post/{POST_ID}/unvote - рейтинг постп вверх
func RegisterVoteHandlers(r *routing.RouteGroup, service vote.IService, postService post.IService, logger log.ILogger, authHandler routing.Handler) {
	c := voteController{
		Service:     service,
		PostService: postService,
		Logger:      logger,
	}

	r.Use(authHandler)

	r.Get(`/post/<postId:\d+>/upvote`, c.upvote)
	r.Get(`/post/<postId:\d+>/downvote`, c.downvote)
	r.Get(`/post/<postId:\d+>/unvote`, c.unvote)
}

func (c voteController) upvote(ctx *routing.Context) error {
	postId, err := c.parseUint(ctx, "postId")
	if err != nil {
		return errorshandler.BadRequest("postId must be uint")
	}
	return c.vote(ctx, postId, 1)
}

func (c voteController) downvote(ctx *routing.Context) error {
	postId, err := c.parseUint(ctx, "postId")
	if err != nil {
		return errorshandler.BadRequest("postId must be uint")
	}
	return c.vote(ctx, postId, -1)
}

func (c voteController) vote(ctx *routing.Context, postId uint, val int) error {
	entity := c.Service.NewEntity(postId, val)
	session := auth.CurrentSession(ctx.Request.Context())
	entity.UserID = session.UserID
	entity.User = session.User

	if err := c.Service.Vote(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError(err.Error())
	}

	post, err := c.PostService.Get(ctx.Request.Context(), postId)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(post, http.StatusOK)
}

func (c voteController) unvote(ctx *routing.Context) error {
	postId, err := c.parseUint(ctx, "postId")
	if err != nil {
		return errorshandler.BadRequest("postId must be uint")
	}
	session := auth.CurrentSession(ctx.Request.Context())
	entity := &vote.Vote{
		PostID: postId,
		UserID: session.UserID,
		User:   session.User,
	}

	if err := c.Service.Unvote(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError(err.Error())
	}

	post, err := c.PostService.Get(ctx.Request.Context(), postId)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(post, http.StatusOK)
}
