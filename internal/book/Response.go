package book

import (
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
)

type bookResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Page        int       `json:"page"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func NewBookResponse(book Book) bookResponse {
	return bookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Page:        book.Page,
		CreatedAt:   book.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   book.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (response bookResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewBookListResponse(books []Book) []render.Renderer {
	list := []render.Renderer{}
	for _, book := range books {
		list = append(list, NewBookResponse(book))
	}

	return list
}
