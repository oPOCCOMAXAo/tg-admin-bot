package tg

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/logger"
	xmodels "github.com/opoccomaxao/tg-admin-bot/pkg/models"
	"github.com/opoccomaxao/tg-admin-bot/pkg/tg/middleware"
	"github.com/opoccomaxao/tg-instrumentation/router"
	"github.com/pkg/errors"
)

type Service struct {
	config Config
	client *bot.Bot
	router *router.Router
	logger *slog.Logger
}

type Config struct {
	Debug     bool   `env:"DEBUG"`
	NoInit    bool   `env:"NO_INIT"`
	Token     string `env:"TOKEN,required"`
	ServerURL string `env:"SERVER_URL"`
	HookURL   string `env:"HOOK_URL"`
}

func New(
	config Config,
	logger *slog.Logger,
) (*Service, error) {
	res := &Service{
		config: config,
		logger: logger.WithGroup("tg"),
	}

	err := res.initClient()
	if err != nil {
		return nil, err
	}

	res.initRouter()

	return res, nil
}

func (s *Service) initClient() error {
	var err error

	opts := []bot.Option{
		bot.WithSkipGetMe(),
		bot.WithDebugHandler(logger.AsPrintf(s.logger.Debug)),
		bot.WithErrorsHandler(s.ErrorHandler),
	}

	if s.config.Debug {
		opts = append(opts, bot.WithDebug())
	}

	if s.config.ServerURL != "" {
		opts = append(opts, bot.WithServerURL(s.config.ServerURL))
	}

	s.client, err = bot.New(s.config.Token, opts...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) initRouter() {
	opts := []router.Option{
		router.WithClient(s.client),
	}

	if s.config.Debug {
		opts = append(opts, router.WithDebug())
	}

	s.router = router.New(opts...)

	s.router.Use(
		router.Recover(),
		router.AutoAccept(),
		middleware.Telemetry(s.logger),
	)
}

func (s *Service) Client() *bot.Bot {
	return s.client
}

func (s *Service) Router() *router.Router {
	return s.router
}

func (s *Service) OnStart(ctx context.Context) error {
	if s.config.NoInit {
		return nil
	}

	err := s.router.UpdateCommandsDescription(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	ok, err := s.client.SetWebhook(ctx, &bot.SetWebhookParams{
		URL:         s.config.HookURL,
		SecretToken: "",
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if !ok {
		return errors.Wrap(xmodels.ErrFailed, "start webhook")
	}

	ok, err = s.client.SetMyDefaultAdministratorRights(ctx, &bot.SetMyDefaultAdministratorRightsParams{
		Rights: &models.ChatAdministratorRights{
			CanRestrictMembers: true,
			CanManageChat:      true,
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if !ok {
		return errors.Wrap(xmodels.ErrFailed, "set default admin rights")
	}

	return nil
}

func (s *Service) ErrorHandler(err error) {
	s.logger.Error("tg error", slog.Any("error", err))
}
