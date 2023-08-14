package todo

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"goexemples/pkg/Beluga"
	"html/template"
	"net/http"
)

type TodoWebRoutes struct {
	controller *TodoController
}

func NewTodoWebRoutes(controller *TodoController) *TodoWebRoutes {
	return &TodoWebRoutes{
		controller: controller,
	}
}

func (webRoutes *TodoWebRoutes) Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/", webRoutes.List)
	router.Post("/", webRoutes.Create)
	router.Post("/toggle", webRoutes.Toggle)

	return router
}

func RenderPage(writer http.ResponseWriter, request *http.Request, templatePath string, data interface{}) {
	defaultLayoutPath := "internal/todo/templates/layouts/default.gohtml"
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	tmpl, err := template.ParseFiles(defaultLayoutPath, templatePath, "internal/todo/templates/pages/todo-list-item.gohtml")
	if err != nil {
		render.Render(writer, request, Beluga.ErrNotFound)
		return
	}

	err = tmpl.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}
}

func RenderListItem(writer http.ResponseWriter, request *http.Request, data interface{}) {
	tmpl, err := template.ParseFiles("internal/todo/templates/pages/todo-list-item.gohtml")
	if err != nil {
		render.Render(writer, request, Beluga.ErrNotFound)
		return
	}

	err = tmpl.ExecuteTemplate(writer, "item", data)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}
}

func (webRoutes *TodoWebRoutes) List(writer http.ResponseWriter, request *http.Request) {
	todos, err := webRoutes.controller.List(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	RenderPage(writer, request, "internal/todo/templates/pages/todo-list.gohtml", todos)
}

func (webRoutes *TodoWebRoutes) Toggle(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	todo, err := webRoutes.controller.Toggle(request.Context(), uuid.MustParse(id))
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}
	RenderListItem(writer, request, todo)
}

func (webRoutes *TodoWebRoutes) Create(writer http.ResponseWriter, request *http.Request) {
	var cerateTodoParams CreateTodoParams
	err := json.NewDecoder(request.Body).Decode(&cerateTodoParams)

	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	todo, err := webRoutes.controller.Create(request.Context(), cerateTodoParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusCreated)
	render.Render(writer, request, NewTodoResponse(todo))
}
