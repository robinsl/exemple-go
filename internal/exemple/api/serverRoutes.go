package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (server *Server) routes() {
	server.router.Use(render.SetContentType(render.ContentTypeJSON))

	server.router.Route("/books", func(router chi.Router) {
		router.Get("/", server.handleAllBooks)
		//router.Post("/", server.handleAddBook)
		//router.Routes("/{id}", func(router chi.Router) {
		//	router.Get("/", server.handleGetBook)
		//	router.Put("/", server.handleUpdateBook)
		//	router.Delete("/", server.handleDeleteBook)
		//}
	})
}
