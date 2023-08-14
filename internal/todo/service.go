package todo

import "goexemples/pkg/Beluga"

type TodoService struct {
	database   Beluga.Database
	store      *TodoStore
	controller *TodoController
	ApiRoutes  *TodoApiRoutes
	WebRoutes  *TodoWebRoutes
}

func NewTodoService() *TodoService {
	databaseConfiguration, err := Beluga.LoadDatabaseConfiguration("todo")
	if err != nil {
		panic(err)
	}
	database := Beluga.NewDatabase(databaseConfiguration)
	store := NewTodoStore(database)
	controller := NewTodoController(store)
	apiRoutes := NewTodoApiRoutes(controller)
	webRoutes := NewTodoWebRoutes(controller)

	return &TodoService{
		database:   database,
		store:      store,
		controller: controller,
		ApiRoutes:  apiRoutes,
		WebRoutes:  webRoutes,
	}
}
