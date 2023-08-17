package todo

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"goexemples/pkg/Beluga"
	"net/http"
	"strconv"
)

type TodoWebRoutes struct {
	controller        *TodoController
	pageListTemplater *Beluga.Templater
	listTemplater     *Beluga.Templater
	listItemTemplater *Beluga.Templater
}

func NewTodoWebRoutes(controller *TodoController) *TodoWebRoutes {
	pageListData := Beluga.PageData{
		PageHeader: Beluga.PageHeader{
			Title: "Todo List",
		},
	}
	pageTemplater := Beluga.
		NewTemplater("default").
		AddComponent("todo/list").
		AddComponent("todo/mainContent").
		SetTemplateName("Layout").
		SetPageData(pageListData).
		Freeze()
	listTemplater := Beluga.
		NewTemplater("default").
		AddComponent("todo/list").
		SetTemplateName("TodoList").
		Freeze()
	listItemTemplater := Beluga.
		NewTemplater("default").
		AddComponent("todo/list").
		SetTemplateName("TodoListItem").
		Freeze()

	return &TodoWebRoutes{
		controller:        controller,
		pageListTemplater: pageTemplater,
		listTemplater:     listTemplater,
		listItemTemplater: listItemTemplater,
	}
}

func (webRoutes *TodoWebRoutes) Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/", webRoutes.ListPage)
	router.Route("/hx", func(router chi.Router) {
		router.Get("/", webRoutes.List)
		router.Post("/toggle", webRoutes.Toggle)
		router.Get("/count-active", webRoutes.CountActive)
	})
	return router
}

func (webRoutes *TodoWebRoutes) ListPage(writer http.ResponseWriter, request *http.Request) {
	todos, err := webRoutes.controller.List(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	count, err := webRoutes.controller.GetAllActive(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	responseData := struct {
		Todos []Todo
		Count int
	}{
		Todos: todos,
		Count: len(count),
	}

	webRoutes.pageListTemplater.Render(writer, request, responseData)
}

func (webRoutes *TodoWebRoutes) List(writer http.ResponseWriter, request *http.Request) {
	todos, err := webRoutes.controller.List(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	responseData := struct {
		Todos []Todo
	}{
		Todos: todos,
	}

	webRoutes.listTemplater.Render(writer, request, responseData)
}

func (webRoutes *TodoWebRoutes) Toggle(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	todo, err := webRoutes.controller.Toggle(request.Context(), uuid.MustParse(id))
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	writer.Header().Set("HX-Trigger", "todo-toggled")
	webRoutes.listItemTemplater.Render(writer, request, todo)
}

func (webRoutes *TodoWebRoutes) CountActive(writer http.ResponseWriter, request *http.Request) {
	activeList, err := webRoutes.controller.GetAllActive(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusOK)
	render.PlainText(writer, request, strconv.Itoa(len(activeList)))
}
