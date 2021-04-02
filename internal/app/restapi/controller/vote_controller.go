package controller

import (
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/vote"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
)

type voteController struct {
	Service     vote.IService
	PostService post.IService
	Logger      log.ILogger
}

func RegisterVoteHandlers(r *routing.RouteGroup, service vote.IService, postService post.IService, logger log.ILogger, authHandler routing.Handler) {
	/*c := voteController{
		Service:     service,
		PostService: postService,
		Logger:      logger,
	}

	r.Use(authHandler)
	*/
}
