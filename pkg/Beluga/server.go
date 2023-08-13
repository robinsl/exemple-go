package Beluga

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	config HttpServerConfiguration
}

func NewServer(config HttpServerConfiguration) *Server {
	server := &Server{
		config: config,
	}

	return server
}

func (server *Server) Serve(context context.Context, routes chi.Router) {
	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%d", server.config.Port),
		IdleTimeout:  server.config.IdleTimeout,
		ReadTimeout:  server.config.ReadTimeout,
		WriteTimeout: server.config.WriteTimeout,
		Handler:      routes,
	}

	log.Printf("Server starting on port %d...", server.config.Port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func handleShutdown(onShutdownSignal func()) <-chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

		<-shutdownSignal

		onShutdownSignal()
		close(shutdown)
	}()

	return shutdown
}
