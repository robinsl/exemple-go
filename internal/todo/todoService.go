package todo

import "goexemples/pkg/Beluga"

type TodoService struct {
	database Beluga.Database
	store    *TodoStore
	Routes   *TodoCrudRoutes
}

func NewTodoService() *TodoService {
	databaseConfiguration, err := Beluga.LoadDatabaseConfiguration("todo")
	if err != nil {
		panic(err)
	}
	database := Beluga.NewDatabase(databaseConfiguration)
	store := NewTodoStore(database)
	routes := NewTodoCrudRoutes(store)

	return &TodoService{
		database: database,
		store:    store,
		Routes:   routes,
	}
}
