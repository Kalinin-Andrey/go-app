package query

import (
	"context"
)

const QueryParams = "QueryParams"

type Params struct {
	Where	map[string]interface{}
	SortOrder	string
}


func ExtractParams(ctx context.Context) *Params {
	var params *Params
	if p, ok := ctx.Value(QueryParams).(Params); ok {
		params = &p
	}
	return params
}

