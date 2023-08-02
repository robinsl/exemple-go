package store

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Book struct {
	ID        uuid.UUID `bson:"_id"`
	Title     string
	Page      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateBookParams struct {
	Title string
	Page  int
}

type UpdateBookParams struct {
	Title string
	Page  int
}

type BookCrud interface {
	Create(ctx context.Context, params CreateBookParams) (Book, error)
	Read(ctx context.Context, id uuid.UUID) (Book, error)
	Update(ctx context.Context, id uuid.UUID, params UpdateBookParams) (Book, error)
	Delete(ctx context.Context, id uuid.UUID) error
	All(ctx context.Context) ([]Book, error)
}
