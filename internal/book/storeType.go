package book

import (
	"github.com/google/uuid"
	"time"
)

type Book struct {
	ID          uuid.UUID `bson:"_id"`
	Title       string
	Description string
	Page        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateBookParams struct {
	Title       string
	Description string
	Page        int
}

type UpdateBookParams struct {
	Title       string
	Description string
	Page        int
}
