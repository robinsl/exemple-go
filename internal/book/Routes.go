package book

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"goexemples/pkg/Beluga"
	"net/http"
)

type BookCrudRoutes struct {
	store *BookStore
}

func NewBookCrudRoutes(store *BookStore) *BookCrudRoutes {
	return &BookCrudRoutes{
		store: store,
	}
}

func (crudRoute *BookCrudRoutes) Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/", crudRoute.List)
	router.Post("/", crudRoute.Create)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(crudRoute.ResourceCtx)
		router.Get("/", crudRoute.Get)
		router.Put("/", crudRoute.Update)
		router.Delete("/", crudRoute.Delete)
	})

	return router
}

func (crudRoute *BookCrudRoutes) ResourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		bookId := chi.URLParam(request, "id")
		book, err := crudRoute.store.Read(request.Context(), uuid.MustParse(bookId))
		if err != nil {
			render.Render(writer, request, Beluga.ErrNotFound)
			return
		}

		ctx := context.WithValue(request.Context(), "book", &book)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (crudRoute *BookCrudRoutes) List(writer http.ResponseWriter, request *http.Request) {
	books, err := crudRoute.store.List(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.RenderList(writer, request, NewBookListResponse(books))
}

func (crudRoute *BookCrudRoutes) Create(writer http.ResponseWriter, request *http.Request) {
	var cerateBookParams CreateBookParams
	err := json.NewDecoder(request.Body).Decode(&cerateBookParams)

	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	book, err := crudRoute.store.Create(request.Context(), cerateBookParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusCreated)
	render.Render(writer, request, NewBookResponse(book))
}

func (crudRoute *BookCrudRoutes) Get(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	book, ok := ctx.Value("book").(*Book)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewBookResponse(*book))

}

func (crudRoute *BookCrudRoutes) Update(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	book, ok := ctx.Value("book").(*Book)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	var updateBookParams UpdateBookParams
	err := json.NewDecoder(request.Body).Decode(&updateBookParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	bookUpdated, err := crudRoute.store.Update(ctx, book.ID, updateBookParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewBookResponse(bookUpdated))
}

func (crudRoute *BookCrudRoutes) Delete(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	book, ok := ctx.Value("book").(*Book)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	err := crudRoute.store.Delete(ctx, book.ID)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusNoContent)
	render.Render(writer, request, NewBookResponse(*book))
}
