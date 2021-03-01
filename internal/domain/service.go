package domain

import (
	"github.com/minipkg/go-app-common/log"
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
