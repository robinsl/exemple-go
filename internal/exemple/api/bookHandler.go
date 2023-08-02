package api

import (
	"github.com/go-chi/render"
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

func (server *Server) handleAllBooks(writer http.ResponseWriter, request *http.Request) {
	books, err := server.store.All(request.Context())
	if err != nil {
		render.Render(writer, request, ErrInternalServerError)
		return
	}

	render.RenderList(writer, request, NewBookListResponse(books))
}
