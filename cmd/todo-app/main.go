package main

import (
	"goexemples/internal/todo"
	"goexemples/pkg/Beluga"
)

func main() {
	beluga := Beluga.NewBeluga()
	beluga.UseDefaultMiddleWare()

	TodoService := todo.NewTodoService()

	beluga.MountRoutes("/", TodoService.ApiRoutes)
	beluga.Serve()
}
