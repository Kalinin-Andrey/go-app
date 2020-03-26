package controller

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"strconv"
)

type IService interface {}

type Controller struct {
}


func (c voteController) parseUint(ctx *routing.Context, paramName string) (uint, error) {
	paramVal, err := strconv.ParseUint(ctx.Param(paramName), 10, 64)
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return 0, err
	}
	return uint(paramVal), nil
}


func (c Controller) ExtractQueryFromContext(ctx *routing.Context) map[string]interface{} {
	query := make(map[string]interface{}, 1)

	for _, paramName := range matchedParams {

		if paramVal := ctx.Param(paramName); paramVal != "" {
			query[paramName] = paramVal
		}
	}

	return query
}


