package todo

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	ID        uuid.UUID `bson:"_id"`
	Title     string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTodoParams struct {
	Title string
}

type UpdateTodoParams struct {
	Title  string
	Status bool
}
