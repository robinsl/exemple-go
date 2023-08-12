package todo

import (
	"github.com/go-chi/chi/v5"
	"log"
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

	router.Get("/", crudRoute.List)
	router.Post("/", crudRoute.Create)

	router.Route("/{todosId}", func(router chi.Router) {
		router.Use(crudRoute.ResourceCtx)
		router.Get("/", crudRoute.Get)
		router.Put("/", crudRoute.Update)
		router.Delete("/", crudRoute.Delete)
	})

	return router
}

func (crudRoute *TodoCrudRoutes) ResourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Print("todosCtx")
		next.ServeHTTP(writer, request)
	})
}

func (crudRoute *TodoCrudRoutes) List(writer http.ResponseWriter, request *http.Request) {
	todos, err := crudRoute.store.List(request.Context())
	if err != nil {
		log.Print(err)
		return
	}

	log.Print(todos)
}

func (crudRoute *TodoCrudRoutes) Create(writer http.ResponseWriter, request *http.Request) {
	log.Print("todos_create")
}

func (crudRoute *TodoCrudRoutes) Get(writer http.ResponseWriter, request *http.Request) {
	log.Print("todos_get")
}

func (crudRoute *TodoCrudRoutes) Update(writer http.ResponseWriter, request *http.Request) {
	log.Print("todos_update")
}

func (crudRoute *TodoCrudRoutes) Delete(writer http.ResponseWriter, request *http.Request) {
	log.Print("todos_delete")
}
