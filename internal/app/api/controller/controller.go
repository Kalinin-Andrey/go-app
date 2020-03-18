package controller

import (
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

type IService interface {}

type Controller struct {
	Service IService
	Logger  log.ILogger
}



