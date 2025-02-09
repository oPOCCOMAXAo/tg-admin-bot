package tg

import (
	"context"
)

func (s *Service) SetupCommands(ctx context.Context) {
	if s.client == nil {
		return
	}

	for _, value := range s.router.ListCommandsParams() {
		_, _ = s.client.SetMyCommands(ctx, value)
	}
}
