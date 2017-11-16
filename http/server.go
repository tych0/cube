package http

import (
	"fmt"
	"net/http"

	"github.com/anuvu/cube/config"
	"github.com/anuvu/cube/service"
)

// Service is the object through which people can register HTTP servers.
type Service interface {
	Register(string, http.Handler)
}

type server struct {
	port int
	mux  *http.ServeMux
	server http.Server
	running bool
}

// NewService creates a new HTTP Service
func NewService(ctx service.Context) Service {
	s := &server{mux: http.NewServeMux()}
	ctx.AddLifecycle(s.getLifecycle())
	return s
}

func (s *server) Register(url string, h http.Handler) {
	s.mux.Handle(url, h)
}

func (s *server) getLifecycle() *service.Lifecycle {
	return &service.Lifecycle {
		ConfigHook: s.ConfigHook,
		StartHook: s.StartHook,
		StopHook: s.StopHook,
		HealthHook: s.HealthHook,
	}
}

func (s *server) ConfigHook(store config.Store) error {
	if err := store.Get("httpPort", &s.port); err != nil {
		return err
	}

	return nil
}

func (s *server) StartHook() error {
	s.server = http.Server{Addr: fmt.Sprintf("localhost:%d", s.port), Handler: s.mux}
	go func() {
		s.running = true
		err := s.server.ListenAndServe()
		if err != nil {
			fmt.Println("serve failed %v", err)
		}
		s.running = false
	}()
	return nil
}

func (s *server) StopHook() error {
	return s.server.Close()
}

func (s *server) HealthHook() bool {
	return s.running
}
