package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"service/pkg"
)

type Server struct {
	host   string
	port   string
	logger *pkg.Logger
	Router *mux.Router
}

func (s *Server) AddMiddleware(mw mux.MiddlewareFunc) {
	s.Router.Use(mw)
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    s.host + ":" + s.port,
		Handler: s.Router,
	}

	go func() {
		s.logger.Info("Running server")

		if err := httpServer.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("Error: %s. Server has been stopped", err))
		}
	}()
}

func CreateNewServer(host, port string) *Server {
	server := &Server{
		host:   host,
		port:   port,
		logger: pkg.GetLogger(),
		Router: mux.NewRouter(),
	}

	server.AddMiddleware(ResponseCheckerMiddleware)

	return server
}
