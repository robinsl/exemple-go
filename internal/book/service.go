package book

import "goexemples/pkg/Beluga"

type BookService struct {
	database Beluga.Database
	store    *BookStore
	Routes   *BookCrudRoutes
}

func NewBookService() *BookService {
	databaseConfiguration, err := Beluga.LoadDatabaseConfiguration("book")
	if err != nil {
		panic(err)
	}
	database := Beluga.NewDatabase(databaseConfiguration)
	store := NewBookStore(database)
	routes := NewBookCrudRoutes(store)

	return &BookService{
		database: database,
		store:    store,
		Routes:   routes,
	}
}
