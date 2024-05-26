package native

import (
	"context"
	"fmt"
	"mr-tasker/api/native/handlers"
	"net/http"
	"strconv"
)

type Server struct {
	port    int
	mux     *http.ServeMux
	server  *http.Server
	handler *handlers.Handler
}

func NewHttpServer(port int) *Server {
	return &Server{
		mux:     http.NewServeMux(),
		port:    port,
		handler: &handlers.Handler{},
	}

}

func (s *Server) addHandler(p string, h http.Handler) {
	s.mux.Handle(p, h)
}

func (s *Server) Serve() error {
	// set up routing
	s.addHandler("/status", http.HandlerFunc(s.handler.StatusHandler()))

	// start server
	p := strconv.Itoa(s.port)
	s.server = &http.Server{Addr: fmt.Sprintf(":%s", p), Handler: s.mux}

	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
