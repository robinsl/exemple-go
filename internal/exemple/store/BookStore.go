package store

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goexemples/internal/exemple/config"
	"log"
	"time"
)

type BookStore struct {
	database   config.Database
	client     *mongo.Client
	collection *mongo.Collection
}

func NewBookStore(database config.Database) *BookStore {
	return &BookStore{
		database: database,
	}
}

func (store *BookStore) Connect(ctx context.Context) error {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(store.database.ConnectionString).SetServerAPIOptions(serverApi))
	if err != nil {
		return err
	}

	store.client = client

	// Ping the database to check if the connection was successful.
	err = store.client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("Failed to ping the database:", err)
		log.Fatal("ConnectionString:", store.database.ConnectionString)
	}

	store.collection = client.Database(store.database.DatabaseName).Collection("books")

	return nil
}

func (store *BookStore) Disconnect(ctx context.Context) error {
	if err := store.client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (store *BookStore) All(ctx context.Context) ([]Book, error) {
	err := store.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Disconnect(ctx)

	cursor, err := store.collection.Find(ctx, bson.D{})
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
	err := store.Connect(ctx)
	if err != nil {
		return Book{}, err
	}
	defer store.Disconnect(ctx)

	bookUuid, err := uuid.NewRandom()
	if err != nil {
		return Book{}, err
	}

	book := Book{
		ID:        bookUuid,
		Title:     params.Title,
		Page:      params.Page,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = store.collection.InsertOne(ctx, book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (store *BookStore) Read(ctx context.Context, id uuid.UUID) (Book, error) {
	err := store.Connect(ctx)
	if err != nil {
		return Book{}, err
	}
	defer store.Disconnect(ctx)

	var book Book
	if err := store.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&book); err != nil {
		return Book{}, err
	}

	return book, nil
}

func (store *BookStore) Update(ctx context.Context, id uuid.UUID, params UpdateBookParams) (Book, error) {
	err := store.Connect(ctx)
	if err != nil {
		return Book{}, err
	}
	defer store.Disconnect(ctx)

	book, err := store.Read(ctx, id)
	if err != nil {
		return Book{}, err
	}

	book.Title = params.Title
	book.Page = params.Page
	book.UpdatedAt = time.Now().UTC()

	_, err = store.collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", book}})
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (store *BookStore) Delete(ctx context.Context, id uuid.UUID) error {
	err := store.Connect(ctx)
	if err != nil {
		return err
	}
	defer store.Disconnect(ctx)

	_, err = store.collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return nil
}
