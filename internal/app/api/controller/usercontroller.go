package controller

import (

	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pkg/errors"
	"strconv"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

type userController struct {
	Service user.IService
	Logger  log.ILogger
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterUserHandlers(r *routing.RouteGroup, service user.IService, logger log.ILogger) {
	c := userController{
		Service:	service,
		Logger:		logger,
	}

	r.Get(`/user/<id:\d+>`, c.get)
	r.Get("/users", c.list)

}

// get method is for a getting a one enmtity by ID
func (c userController) get(context *routing.Context) error {
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		return errors.Wrapf(err, "Can not parse uint64 from %q", context.Param("id"))
	}
	entity, err := c.Service.Get(context.Request.Context(), uint(id))
	if err != nil {
		return err
	}
	context.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return context.Write(entity)
}

// list method is for a getting a list of all entities
func (c userController) list(context *routing.Context) error {
	ctx := context.Request.Context()
	items, err := c.Service.List(ctx)
	if err != nil {
		return err
	}
	context.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return context.Write(items)
}

