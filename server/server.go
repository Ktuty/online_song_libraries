package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	HttpServer *http.Server
}

func (server *Server) Run(port string, handler http.Handler) error {
	server.HttpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return server.HttpServer.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.HttpServer.Shutdown(ctx)
}
