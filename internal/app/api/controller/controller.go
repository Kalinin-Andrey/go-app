package controller

import routing "github.com/go-ozzo/ozzo-routing/v2"

type IService interface {}

type Controller struct {
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


