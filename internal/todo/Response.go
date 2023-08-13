package todo

import (
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
)

type todoResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Status    bool      `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func NewTodoResponse(todo Todo) todoResponse {
	return todoResponse{
		ID:        todo.ID,
		Title:     todo.Title,
		Status:    todo.Status,
		CreatedAt: todo.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: todo.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (response todoResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewTodoListResponse(todos []Todo) []render.Renderer {
	list := []render.Renderer{}
	for _, todo := range todos {
		list = append(list, NewTodoResponse(todo))
	}

	return list
}
