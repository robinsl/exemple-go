package todo

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"goexemples/pkg/Beluga"
	"net/http"
)

type TodoCrudRoutes struct {
	store *TodoStore
}

func NewTodoCrudRoutes(store *TodoStore) *TodoCrudRoutes {
	return &TodoCrudRoutes{
		store: store,
	}
}

func (crudRoute *TodoCrudRoutes) Routes() chi.Router {
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

func (crudRoute *TodoCrudRoutes) ResourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		todoId := chi.URLParam(request, "id")
		todo, err := crudRoute.store.Read(request.Context(), uuid.MustParse(todoId))
		if err != nil {
			render.Render(writer, request, Beluga.ErrNotFound)
			return
		}

		ctx := context.WithValue(request.Context(), "todo", &todo)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (crudRoute *TodoCrudRoutes) List(writer http.ResponseWriter, request *http.Request) {
	todos, err := crudRoute.store.List(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.RenderList(writer, request, NewTodoListResponse(todos))
}

func (crudRoute *TodoCrudRoutes) Create(writer http.ResponseWriter, request *http.Request) {
	var cerateTodoParams CreateTodoParams
	err := json.NewDecoder(request.Body).Decode(&cerateTodoParams)

	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	todo, err := crudRoute.store.Create(request.Context(), cerateTodoParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusCreated)
	render.Render(writer, request, NewTodoResponse(todo))
}

func (crudRoute *TodoCrudRoutes) Get(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	todo, ok := ctx.Value("todo").(*Todo)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewTodoResponse(*todo))

}

func (crudRoute *TodoCrudRoutes) Update(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	todo, ok := ctx.Value("todo").(*Todo)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	var updateTodoParams UpdateTodoParams
	err := json.NewDecoder(request.Body).Decode(&updateTodoParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	todoUpdated, err := crudRoute.store.Update(ctx, todo.ID, updateTodoParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewTodoResponse(todoUpdated))
}

func (crudRoute *TodoCrudRoutes) Delete(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	todo, ok := ctx.Value("todo").(*Todo)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	err := crudRoute.store.Delete(ctx, todo.ID)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusNoContent)
	render.Render(writer, request, NewTodoResponse(*todo))
}
