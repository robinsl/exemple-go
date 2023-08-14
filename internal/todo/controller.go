package todo

import (
	"context"
	"github.com/google/uuid"
)

type TodoController struct {
	store *TodoStore
}

func NewTodoController(store *TodoStore) *TodoController {
	return &TodoController{
		store: store,
	}
}

func (controller *TodoController) List(context context.Context) ([]Todo, error) {
	return controller.store.List(context)
}

func (controller *TodoController) Create(context context.Context, params CreateTodoParams) (Todo, error) {
	return controller.store.Create(context, params)
}

func (controller *TodoController) Get(context context.Context, id uuid.UUID) (Todo, error) {
	return controller.store.Read(context, id)
}

func (controller *TodoController) Update(context context.Context, id uuid.UUID, params UpdateTodoParams) (Todo, error) {
	return controller.store.Update(context, id, params)
}
func (controller *TodoController) Toggle(context context.Context, id uuid.UUID) (Todo, error) {
	return controller.store.Toggle(context, id)
}

func (controller *TodoController) Delete(context context.Context, id uuid.UUID) error {
	return controller.store.Delete(context, id)
}
