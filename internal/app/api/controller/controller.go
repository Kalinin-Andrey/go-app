package controller

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"redditclone/internal/domain"
	"strconv"
)

type IService interface{}

type Controller struct {
}

var matchedParams = []string{
	"userName",
	"category",
}

func (c voteController) parseUint(ctx *routing.Context, paramName string) (uint, error) {
	paramVal, err := strconv.ParseUint(ctx.Param(paramName), 10, 64)
	if err != nil {
		c.Logger.With(ctx.Request.Context()).Info(err)
		return 0, err
	}
	return uint(paramVal), nil
}

func (c Controller) ExtractQueryFromRoutingContext(ctx *routing.Context) domain.DBQueryConditions {
	query := map[string]interface{}{}

	for _, paramName := range matchedParams {

		if paramVal := ctx.Param(paramName); paramVal != "" {
			query[paramName] = paramVal
		}
	}

	return domain.DBQueryConditions{
		Where: query,
	}
}
