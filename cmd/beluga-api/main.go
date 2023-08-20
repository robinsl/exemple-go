package main

import (
	"goexemples/internal/book"
	"goexemples/internal/todo"
	"goexemples/pkg/Beluga"
)

func main() {
	beluga := Beluga.NewBeluga()
	beluga.UseDefaultMiddleWare()

	TodoService := todo.NewTodoService()
	BookService := book.NewBookService()

	beluga.MountRoutes("/todos", TodoService.ApiRoutes)
	beluga.MountRoutes("/todos-app", TodoService.WebRoutes)
	beluga.MountRoutes("/books", BookService.Routes)

	beluga.MountStatic()
	beluga.Serve()
}
