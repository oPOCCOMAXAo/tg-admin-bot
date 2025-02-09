package middleware

import (
	"log/slog"
	"time"

	"github.com/opoccomaxao/tg-instrumentation/router"
)

func Telemetry(
	logger *slog.Logger,
) router.Handler {
	return func(ctx *router.Context) {
		update := ctx.Update()

		args := []any{
			slog.String("pattern", ctx.Pattern()),
		}

		{
			raw := ctx.RawDebug()
			if len(raw) > 0 {
				args = append(args, slog.String("raw", string(raw)))
			}
		}

		switch {
		case update.Message != nil:
			args = append(args,
				slog.String("update_type", "message"),
				slog.Int64("user_id", update.Message.From.ID),
				slog.String("user_name", update.Message.From.Username),
				slog.String("message_text", update.Message.Text),
				slog.Time("message_date", time.Unix(int64(update.Message.Date), 0)),
			)
		case update.CallbackQuery != nil:
			args = append(args,
				slog.String("update_type", "callback"),
				slog.Int64("user_id", update.CallbackQuery.From.ID),
				slog.String("user_name", update.CallbackQuery.From.Username),
				slog.String("data", update.CallbackQuery.Data),
			)
		case update.InlineQuery != nil:
			args = append(args,
				slog.String("update_type", "inline"),
				slog.Int64("user_id", update.InlineQuery.From.ID),
				slog.String("user_name", update.InlineQuery.From.Username),
				slog.String("query", update.InlineQuery.Query),
			)
		}

		logger.InfoContext(ctx.Context(), "request", args...)

		ctx.Next()
	}
}
