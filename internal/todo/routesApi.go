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

type TodoApiRoutes struct {
	controller *TodoController
}

func NewTodoApiRoutes(controller *TodoController) *TodoApiRoutes {
	return &TodoApiRoutes{
		controller: controller,
	}
}

func (apiRoutes *TodoApiRoutes) Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/", apiRoutes.List)
	router.Post("/", apiRoutes.Create)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(apiRoutes.ResourceCtx)
		router.Get("/", apiRoutes.Get)
		router.Put("/", apiRoutes.Update)
		router.Delete("/", apiRoutes.Delete)
	})

	return router
}

func (apiRoutes *TodoApiRoutes) ResourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		todoId := chi.URLParam(request, "id")
		todo, err := apiRoutes.controller.Get(request.Context(), uuid.MustParse(todoId))
		if err != nil {
			render.Render(writer, request, Beluga.ErrNotFound)
			return
		}

		ctx := context.WithValue(request.Context(), "todo", &todo)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (apiRoutes *TodoApiRoutes) List(writer http.ResponseWriter, request *http.Request) {
	todos, err := apiRoutes.controller.List(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.RenderList(writer, request, NewTodoListResponse(todos))
}

func (apiRoutes *TodoApiRoutes) Create(writer http.ResponseWriter, request *http.Request) {
	var createTodoParams CreateTodoParams
	err := json.NewDecoder(request.Body).Decode(&createTodoParams)

	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	todo, err := apiRoutes.controller.Create(request.Context(), createTodoParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusCreated)
	render.Render(writer, request, NewTodoResponse(todo))
}

func (apiRoutes *TodoApiRoutes) Get(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	todo, ok := ctx.Value("todo").(*Todo)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewTodoResponse(*todo))

}

func (apiRoutes *TodoApiRoutes) Update(writer http.ResponseWriter, request *http.Request) {
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

	todoUpdated, err := apiRoutes.controller.Update(ctx, todo.ID, updateTodoParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Render(writer, request, NewTodoResponse(todoUpdated))
}

func (apiRoutes *TodoApiRoutes) Delete(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	todo, ok := ctx.Value("todo").(*Todo)
	if !ok {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	err := apiRoutes.controller.Delete(ctx, todo.ID)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusNoContent)
	render.Render(writer, request, NewTodoResponse(*todo))
}
