package todo

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"goexemples/pkg/Beluga"
	"log"
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
		router.Post("/", webRoutes.Create)
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

func (webRoutes *TodoWebRoutes) Create(writer http.ResponseWriter, request *http.Request) {
	var createTodoParams CreateTodoParams
	err := request.ParseForm()
	if err != nil {
		log.Println(err)
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	createTodoParams.Title = request.FormValue("title")

	todo, err := webRoutes.controller.Create(request.Context(), createTodoParams)
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	writer.Header().Set("HX-Trigger", "todo-created")
	webRoutes.listItemTemplater.Render(writer, request, todo)
}

func (webRoutes *TodoWebRoutes) Toggle(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")
	todo, err := webRoutes.controller.Toggle(request.Context(), uuid.MustParse(id))
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	writer.Header().Set("HX-Trigger", "todo-toggled")
	http.SetCookie(writer, &http.Cookie{
		Name:  "todo-toggled",
		Path:  "/",
		Value: id,
	})
	webRoutes.listItemTemplater.Render(writer, request, todo)
}

func (webRoutes *TodoWebRoutes) CountActive(writer http.ResponseWriter, request *http.Request) {
	c, err := request.Cookie("todo-toggled")
	if err != nil {
		log.Println(err)
	}
	log.Printf("Loading %s from cookie, Value is: %s", c.Name, c.Value)

	activeList, err := webRoutes.controller.GetAllActive(request.Context())
	if err != nil {
		render.Render(writer, request, Beluga.ErrInternalServerError)
		return
	}

	render.Status(request, http.StatusOK)
	render.PlainText(writer, request, strconv.Itoa(len(activeList)))
}
