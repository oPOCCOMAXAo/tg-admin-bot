package handlers

import (
	"github.com/opoccomaxao/tg-instrumentation/router"
)

func (s *Service) GroupMessage(ctx *router.Context) {
	err := s.domain.ProcessMessageCheck(ctx.Context(), ctx.Update())
	if err != nil {
		ctx.Error(err)

		return
	}
}
