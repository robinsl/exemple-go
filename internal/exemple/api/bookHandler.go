package api

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"goexemples/internal/exemple/store"
	"net/http"
)

type bookResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Page      int    `json:"page"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewBookResponse(book store.Book) bookResponse {
	return bookResponse{
		ID:        book.ID.String(),
		Title:     book.Title,
		Page:      book.Page,
		CreatedAt: book.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: book.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (response bookResponse) Render(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func NewBookListResponse(books []store.Book) []render.Renderer {
	list := []render.Renderer{}
	for _, book := range books {
		list = append(list, NewBookResponse(book))
	}

	return list
}

func (server *Server) bookCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		bookId := chi.URLParam(request, "id")
		book, err := server.store.Read(request.Context(), uuid.MustParse(bookId))
		if err != nil {
			render.Render(writer, request, ErrInternalServerError)
			return
		}

		ctx := context.WithValue(request.Context(), "book", &book)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (server *Server) handleAllBooks(writer http.ResponseWriter, request *http.Request) {
	books, err := server.store.All(request.Context())
	if err != nil {
		render.Render(writer, request, ErrInternalServerError)
		return
	}

	render.RenderList(writer, request, NewBookListResponse(books))
}

func (server *Server) handleGetBook(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	book, ok := ctx.Value("book").(*store.Book)
	if !ok {
		render.Render(writer, request, ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewBookResponse(*book))
}

func (server *Server) handleAddBook(writer http.ResponseWriter, request *http.Request) {
	var createBookParam store.CreateBookParams
	json.NewDecoder(request.Body).Decode(&createBookParam)

	book, err := server.store.Create(request.Context(), createBookParam)
	if err != nil {
		render.Render(writer, request, ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusCreated)
	render.Render(writer, request, NewBookResponse(book))
}

func (server *Server) handleUpdateBook(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	book, ok := ctx.Value("book").(*store.Book)
	if !ok {
		render.Render(writer, request, ErrInternalServerError)
		return
	}

	var updateBookParam store.UpdateBookParams
	json.NewDecoder(request.Body).Decode(&updateBookParam)

	bookUpdated, err := server.store.Update(request.Context(), book.ID, updateBookParam)
	if err != nil {
		render.Render(writer, request, ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewBookResponse(bookUpdated))
}

func (server *Server) handleDeleteBook(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	book, ok := ctx.Value("book").(*store.Book)
	if !ok {
		render.Render(writer, request, ErrInternalServerError)
		return
	}

	err := server.store.Delete(request.Context(), book.ID)
	if err != nil {
		render.Render(writer, request, ErrInternalServerError)
		return
	}
	render.Render(writer, request, SuccessResourceDeleted)
}
