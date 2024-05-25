package native

import (
	"fmt"
	"mr-tasker/api/native/handlers"
	"net/http"
	"strconv"
)

type Server struct {
	port    int
	mux     *http.ServeMux
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

func (s *Server) Serve() {
	s.addHandler("/status", http.HandlerFunc(s.handler.StatusHandler()))
	p := strconv.Itoa(s.port)
	http.ListenAndServe(fmt.Sprintf(":%s", p), s.mux)
}
