package handlers

import "github.com/opoccomaxao/tg-instrumentation/router"

func (s *Service) DeleteSelf(ctx *router.Context) {
	_, _ = ctx.DeleteMessageFromCallback()
}
