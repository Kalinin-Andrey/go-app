package domain

import (
	"redditclone/internal/pkg/log"
)

const MaxLIstLimit = 1000

type Service struct {
	logger log.ILogger
}

type DBQueryConditions struct {
	Where     map[string]interface{}
	SortOrder map[string]interface{}
	Limit     int
	Offset    int
}
