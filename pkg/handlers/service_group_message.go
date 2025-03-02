package handlers

import (
	"github.com/opoccomaxao/tg-instrumentation/router"
)

func (s *Service) GroupMessage(ctx *router.Context) {
	err := s.domain.HandleMessage(ctx.Context(), ctx.Update())
	if err != nil {
		ctx.Error(err)

		return
	}
}
