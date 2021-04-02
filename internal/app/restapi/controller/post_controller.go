package controller

import (
	"net/http"

	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/minipkg/ozzo_routing"

	"redditclone/internal/pkg/apperror"
	"redditclone/internal/pkg/auth"
	"redditclone/internal/pkg/errorshandler"

	"redditclone/internal/domain/post"
	"redditclone/internal/domain/user"
	"redditclone/internal/domain/vote"
)

type postController struct {
	Service     post.IService
	UserService user.IService
	Logger      log.ILogger
}

// RegisterHandlers sets up the routing of the HTTP handlers.
//	GET /api/posts/ - список всех постов
//	GET /api/post/{POST_ID} - детали поста с комментами
//	GET /api/posts/{CATEGORY_NAME} - список постов конкретной категории
//	GET /api/user/{USER_LOGIN} - получение всех постов конкртеного пользователя
//	POST /api/posts/ - добавление поста - обратите внимание - есть с урлом, а есть с текстом
//	DELETE /api/post/{POST_ID} - удаление поста
//	GET /api/post/{POST_ID}/upvote - рейтинг поста вверх
//	GET /api/post/{POST_ID}/downvote - рейтинг поста вниз
//	GET /api/post/{POST_ID}/unvote - рейтинг постп вверх
func RegisterPostHandlers(r *routing.RouteGroup, service post.IService, userService user.IService, logger log.ILogger, authHandler routing.Handler) {
	c := postController{
		Service:     service,
		UserService: userService,
		Logger:      logger,
	}

	r.Get("/posts", c.list)
	r.Get(`/post/<id>`, c.get)
	r.Get(`/posts/<category:\w+>`, c.list)
	r.Get(`/user/<userName:\w+>`, c.list)

	r.Use(authHandler)

	r.Post("/posts", c.create)
	r.Delete(`/post/<id>`, c.delete)

	r.Get(`/post/<postId>/upvote`, c.upvote)
	r.Get(`/post/<postId>/downvote`, c.downvote)
	r.Get(`/post/<postId>/unvote`, c.unvote)
}

// get method is for getting a one entity by ID
func (c *postController) get(ctx *routing.Context) error {
	id := ctx.Param("id")

	entity, err := c.Service.Get(ctx.Request.Context(), id)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	if err = c.Service.ViewsIncr(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.Write(entity)
}

// list method is for a getting a list of all entities
func (c *postController) list(ctx *routing.Context) error {
	rctx := ctx.Request.Context()
	where := c.Service.NewEntity()

	if len(ctx.Request.URL.Query()) > 0 {
		err := ozzo_routing.ParseQueryParamsIntoStruct(ctx, where)
		if err != nil {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.BadRequest("")
		}
	}

	if category := ctx.Param("category"); category != "" {
		where.Category = category
	}

	if userName := ctx.Param("userName"); userName != "" {
		// @todo: make case insensitive
		user, err := c.UserService.First(rctx, &user.User{
			Name: userName,
		})
		if err != nil {
			if err == apperror.ErrNotFound {
				c.Logger.With(ctx.Request.Context()).Info(errors.Wrapf(err, "Can not find user with name: %q", userName))
				return errorshandler.NotFound("Can not find user")
			}
			c.Logger.With(ctx.Request.Context()).Error(err)
			return errorshandler.InternalServerError("")
		}
		where.UserID = user.ID
	}
	cond := selection_condition.SelectionCondition{
		Where: where,
	}

	items, err := c.Service.Query(rctx, cond)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}
	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.Write(items)
}

func (c *postController) create(ctx *routing.Context) error {
	entity := c.Service.NewEntity()
	if err := ctx.Read(entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	if err := entity.Validate(); err != nil {
		return errorshandler.BadRequest(err.Error())
	}

	session := auth.CurrentSession(ctx.Request.Context())
	entity.UserID = session.UserID
	entity.User = session.User

	if err := c.Service.Create(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(entity, http.StatusCreated)
}

func (c *postController) delete(ctx *routing.Context) error {
	id := ctx.Param("id")

	if err := c.Service.Delete(ctx.Request.Context(), id); err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.With(ctx.Request.Context()).Info(err)
			return errorshandler.NotFound("")
		}
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError("")
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(errorshandler.SuccessMessage(), http.StatusOK)
}

func (c *postController) upvote(ctx *routing.Context) error {
	return c.vote(ctx, ctx.Param("postId"), 1)
}

func (c *postController) downvote(ctx *routing.Context) error {
	return c.vote(ctx, ctx.Param("postId"), -1)
}

func (c *postController) vote(ctx *routing.Context, postId string, val int) error {
	session := auth.CurrentSession(ctx.Request.Context())
	entity := c.Service.NewVoteEntity(session.UserID, postId, val)
	entity.User = session.User

	if err := c.Service.Vote(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Error(err)
		return errorshandler.InternalServerError(err.Error())
	}

	post, err := c.Service.Get(ctx.Request.Context(), postId)
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

func (c *postController) unvote(ctx *routing.Context) error {
	postId := ctx.Param("postId")
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

	post, err := c.Service.Get(ctx.Request.Context(), postId)
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
