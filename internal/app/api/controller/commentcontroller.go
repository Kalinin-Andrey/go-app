package controller

import (
	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/auth"
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"
	"net/http"
	"strconv"

	"github.com/go-ozzo/ozzo-routing/v2"

	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

type commentController struct {
	Controller
	Service     comment.IService
	PostService post.IService
	Logger      log.ILogger
}

//	POST /api/post/{POST_ID} - добавление коммента
//	DELETE /api/post/{POST_ID}/{COMMENT_ID} - удаление коммента
func RegisterCommentHandlers(r *routing.RouteGroup, service comment.IService, postService post.IService, logger log.ILogger, authHandler routing.Handler) {
	c := commentController{
		Service:     service,
		PostService: postService,
		Logger:      logger,
	}

	r.Use(authHandler)

	r.Post(`/post/<postId:\d+>`, c.create)
	r.Delete(`/post/<postId:\d+>/<id:\d+>`, c.delete)
}

func (c commentController) create(ctx *routing.Context) error {
	pId, err := strconv.ParseUint(ctx.Param("postId"), 10, 64)
	postId := uint(pId)
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest("postId must be uint")
	}

	entity := c.Service.NewEntity()
	if err := ctx.Read(entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	if err := entity.Validate(); err != nil {
		return errorshandler.BadRequest(err.Error())
	}

	session := auth.CurrentSession(ctx.Request.Context())
	entity.PostID = postId
	entity.UserID = session.UserID
	entity.User = session.User

	if err := c.Service.Create(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
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
	return ctx.WriteWithStatus(post, http.StatusCreated)
}

func (c commentController) delete(ctx *routing.Context) error {
	postId, err := strconv.ParseUint(ctx.Param("postId"), 10, 64)
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest("postId must be uint")
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest("id must be uint")
	}

	if err := c.Service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	post, err := c.PostService.Get(ctx.Request.Context(), uint(postId))
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
