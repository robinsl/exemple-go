package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"goexemples/internal/exemple/config"
	"goexemples/internal/exemple/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	config config.HttpServer
	store  store.BookCrud
	router *chi.Mux
}

func NewServer(config config.HttpServer, store store.BookCrud) *Server {
	server := &Server{
		config: config,
		store:  store,
		router: chi.NewRouter(),
	}

	server.routes()

	return server
}

func (server *Server) Serve(context context.Context) {
	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%d", server.config.Port),
		Handler:      server.router,
		IdleTimeout:  server.config.IdleTimeout,
		ReadTimeout:  server.config.ReadTimeout,
		WriteTimeout: server.config.WriteTimeout,
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

	log.Println("Shutdown gracefully")
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
