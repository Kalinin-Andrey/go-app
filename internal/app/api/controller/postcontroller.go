package controller

import (
	"context"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/pkg/auth"
	"github.com/Kalinin-Andrey/redditclone/pkg/errorshandler"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/go-ozzo/ozzo-routing/v2"

	"github.com/Kalinin-Andrey/redditclone/pkg/log"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
)

type postController struct {
	Controller
	Service post.IService
	UserService user.IService
	Logger  log.ILogger
}

// RegisterHandlers sets up the routing of the HTTP handlers.
//	GET /api/posts/ - список всех постов
//	GET /api/post/{POST_ID} - детали поста с комментами
//	GET /api/posts/{CATEGORY_NAME} - список постов конкретной категории
//	GET /api/user/{USER_LOGIN} - получение всех постов конкртеного пользователя
//	POST /api/posts/ - добавление поста - обратите внимание - есть с урлом, а есть с текстом
//	DELETE /api/post/{POST_ID} - удаление поста

//	POST /api/post/{POST_ID} - добавление коммента
//	DELETE /api/post/{POST_ID}/{COMMENT_ID} - удаление коммента

//	GET /api/post/{POST_ID}/upvote - рейтинг постп вверх
//	GET /api/post/{POST_ID}/downvote - рейтинг поста вниз
func RegisterPostHandlers(r *routing.RouteGroup, service post.IService, userService user.IService, logger log.ILogger, authHandler routing.Handler) {
	c := postController{
		Service:		service,
		UserService:	userService,
		Logger:			logger,
	}

	r.Get("/posts", c.list)
	r.Get(`/posts/<id:\d+>`, c.get)
	r.Get(`/posts/<category:\w+>`, c.list)
	r.Get(`/user/<userName:\w+>`, c.list)

	r.Use(authHandler)

	r.Post("/posts", c.create)
	r.Delete(`/posts/<id:\d+>`, c.delete)
}

var matchedParams = []string{
	"userName",
	"category",
}

// list method is for a getting a list of all posts of exactly user
/*func (c postController) getByCategory(ctx *routing.Context) error {
	rctx := ctx.Request.Context()
	items, err := c.Service.List(rctx)
	if err != nil {
		return err
	}
	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.Write(items)
}*/


// get method is for a getting a one enmtity by ID
func (c postController) get(ctx *routing.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return errors.Wrapf(err, "Can not parse uint64 from %q", ctx.Param("id"))
	}
	entity, err := c.Service.Get(ctx.Request.Context(), uint(id))
	if err != nil {
		return err
	}

	if err = c.Service.ViewsIncr(ctx.Request.Context(), entity); err != nil {
		return err
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.Write(entity)
}

// list method is for a getting a list of all entities
func (c postController) list(ctx *routing.Context) error {
	where := c.ExtractQueryFromContext(ctx)
	rctx := ctx.Request.Context()

	for key, val := range where {
		if key == "userName" {
			// @todo: make the case insensitive search
			userName, ok := val.(string)
			if !ok {
				return errors.Errorf("Can not assert interface{} to string for value: %#v", val)
			}
			user, err := c.UserService.First(rctx, &user.User{
				Name:	userName,
			})
			if err != nil {
				return errors.Wrapf(err, "Can not find user with name: %q", userName)
			}
			delete(where, "userName")
			where["UserID"] = user.ID
		}
	}

	if len(where) == 0 {
		rctx = context.WithValue(rctx, "Where", where)
	}

	items, err := c.Service.List(rctx)
	if err != nil {
		return err
	}
	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.Write(items)
}

func (c postController) create(ctx *routing.Context) error {
	entity := c.Service.NewEntity()
	if err := ctx.Read(entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}
	sessRepo := auth.CurrentSession(ctx.Request.Context())
	entity.UserID	= sessRepo.Session.UserID
	entity.User		= sessRepo.Session.User

	if err := c.Service.Create(ctx.Request.Context(), entity); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest(err.Error())
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus(entity, http.StatusCreated)
}


func (c postController) delete(ctx *routing.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.BadRequest("id must be uint")
	}

	if err := c.Service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return errorshandler.NotFound("")
	}

	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return ctx.WriteWithStatus("OK", http.StatusOK)
}

