package controller

import (
	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
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
	Service comment.IService
	Logger  log.ILogger
}

//	POST /api/post/{POST_ID} - добавление коммента
//	DELETE /api/post/{POST_ID}/{COMMENT_ID} - удаление коммента
func RegisterCommentHandlers(r *routing.RouteGroup, service comment.IService, logger log.ILogger, authHandler routing.Handler) {
	c := commentController{
		Service:		service,
		Logger:			logger,
	}

	r.Use(authHandler)

	r.Post(`/post/<postId:\d+>`, c.create)
	r.Delete(`/post/<postId:\d+>/<id:\d+>`, c.delete)
}


func (c commentController) create(ctx *routing.Context) error {
	postId, err := strconv.ParseUint(ctx.Param("postId"), 10, 64)
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

	sessRepo := auth.CurrentSession(ctx.Request.Context())
	entity.PostID	= uint(postId)
	entity.UserID	= sessRepo.Session.UserID
	entity.User		= sessRepo.Session.User

	if err := c.Service.Create(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(entity, http.StatusCreated)
}


func (c commentController) delete(ctx *routing.Context) error {
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

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus("OK", http.StatusOK)
}


