package server

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Server struct {
	mu sync.Mutex

	config   Config
	server   http.Server
	router   *http.ServeMux
	listener net.Listener
}

type Config struct {
	Port int `env:"PORT" envdefault:"8080"`
}

func New(
	config Config,
) *Server {
	res := Server{
		config: config,
		server: http.Server{
			ReadTimeout: time.Minute,
		},
		router: http.NewServeMux(),
	}

	res.server.Handler = res.router

	return &res
}

func (s *Server) OnStart(
	_ context.Context,
	cancelCauseFn context.CancelCauseFunc,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		return errors.WithStack(ErrServerAlreadyStarted)
	}

	var err error

	s.listener, err = net.ListenTCP("tcp", &net.TCPAddr{
		Port: s.config.Port,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	go func() {
		err := s.server.Serve(s.listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			err = errors.WithStack(err)
		} else {
			err = nil
		}

		cancelCauseFn(err)
	}()

	return nil
}

func (s *Server) OnStop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
