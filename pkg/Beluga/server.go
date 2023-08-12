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
		Addr:    ":3333",
		Handler: routes,
	}

	shutdownCompleted := handleShutdown(func() {
		if err := httpServer.Shutdown(context); err != nil {
			fmt.Printf("Error on shutdown http server: %v", err)
		}
	})

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		<-shutdownCompleted
	} else {
		log.Printf("http.ListenAndServe: %v\n", err)
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
