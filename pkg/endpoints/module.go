package endpoints

import (
	"net/http"

	"github.com/opoccomaxao/tg-instrumentation/router"
	"go.uber.org/fx"
)

func Invoke() fx.Option {
	return fx.Module("endpoints",
		fx.Provide(NewService, fx.Private),
		fx.Invoke(RegisterEndpoints),
	)
}

func RegisterEndpoints(
	router *http.ServeMux,
	tgrouter *router.Router,
	service *Service,
) error {
	router.HandleFunc("GET /health", service.Health)

	router.HandleFunc("POST /webhook", tgrouter.HandlerFunc)

	return nil
}
