package Beluga

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
)

type BelugaHandler interface {
	Routes() chi.Router
}

type beluga struct {
	ctx    context.Context
	server *Server
	Route  chi.Router
}

func NewBeluga() *beluga {
	if err := LoadDotEnv(); err != nil {
		log.Fatal(err)
	}

	httpServerConfiguration, err := LoadHttpServerConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	return &beluga{
		ctx:    context.Background(),
		server: NewServer(httpServerConfiguration),
		Route:  chi.NewRouter(),
	}
}

func (b *beluga) UseDefaultMiddleWare() {
	b.Route.Use(middleware.RequestID)
	b.Route.Use(middleware.RealIP)
	b.Route.Use(middleware.Logger)
	b.Route.Use(middleware.Recoverer)
}

func (b beluga) MountRoutes(path string, handler BelugaHandler) {
	b.Route.Mount(path, handler.Routes())
}

func (b *beluga) Serve() {
	b.server.Serve(b.ctx, b.Route)
}
