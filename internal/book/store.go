package book

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"goexemples/pkg/Beluga"
	"time"
)

type BookStore struct {
	database Beluga.Database
}

func NewBookStore(database Beluga.Database) *BookStore {
	return &BookStore{
		database: database,
	}
}

func (store *BookStore) List(ctx context.Context) ([]Book, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer store.database.Disconnect(ctx)

	cursor, err := store.database.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	return books, nil
}

func (store *BookStore) Create(ctx context.Context, params CreateBookParams) (Book, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return Book{}, err
	}
	defer store.database.Disconnect(ctx)

	bookUuid, err := uuid.NewRandom()
	if err != nil {
		return Book{}, err
	}

	book := Book{
		ID:          bookUuid,
		Title:       params.Title,
		Description: params.Description,
		Page:        params.Page,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	_, err = store.database.Collection.InsertOne(ctx, book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (store *BookStore) Read(ctx context.Context, id uuid.UUID) (Book, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return Book{}, err
	}
	defer store.database.Disconnect(ctx)

	var book Book
	err = store.database.Collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (store *BookStore) Update(ctx context.Context, id uuid.UUID, params UpdateBookParams) (Book, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return Book{}, err
	}
	defer store.database.Disconnect(ctx)

	var book Book
	err = store.database.Collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&book)
	if err != nil {
		return Book{}, err
	}

	book.Title = params.Title
	book.Description = params.Description
	book.Page = params.Page
	book.UpdatedAt = time.Now().UTC()

	_, err = store.database.Collection.ReplaceOne(ctx, bson.D{{"_id", id}}, book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (store *BookStore) Delete(ctx context.Context, id uuid.UUID) error {
	err := store.database.Connect(ctx)
	if err != nil {
		return err
	}
	defer store.database.Disconnect(ctx)

	_, err = store.database.Collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return nil
}
