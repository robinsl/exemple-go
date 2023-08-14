package todo

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"goexemples/pkg/Beluga"
	"time"
)

type TodoStore struct {
	database Beluga.Database
}

func NewTodoStore(database Beluga.Database) *TodoStore {
	return &TodoStore{
		database: database,
	}
}

func (store *TodoStore) List(ctx context.Context) ([]Todo, error) {
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

	var todos []Todo
	if err = cursor.All(ctx, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func (store *TodoStore) Create(ctx context.Context, params CreateTodoParams) (Todo, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return Todo{}, err
	}
	defer store.database.Disconnect(ctx)

	todoUuid, err := uuid.NewRandom()
	if err != nil {
		return Todo{}, err
	}

	todo := Todo{
		ID:        todoUuid,
		Title:     params.Title,
		Status:    false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = store.database.Collection.InsertOne(ctx, todo)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (store *TodoStore) Read(ctx context.Context, id uuid.UUID) (Todo, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return Todo{}, err
	}
	defer store.database.Disconnect(ctx)

	var todo Todo
	err = store.database.Collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&todo)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (store *TodoStore) Update(ctx context.Context, id uuid.UUID, params UpdateTodoParams) (Todo, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return Todo{}, err
	}
	defer store.database.Disconnect(ctx)

	var todo Todo
	err = store.database.Collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&todo)
	if err != nil {
		return Todo{}, err
	}

	todo.Title = params.Title
	todo.Status = params.Status
	todo.UpdatedAt = time.Now().UTC()

	_, err = store.database.Collection.ReplaceOne(ctx, bson.D{{"_id", id}}, todo)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (store *TodoStore) Toggle(ctx context.Context, id uuid.UUID) (Todo, error) {
	err := store.database.Connect(ctx)
	if err != nil {
		return Todo{}, err
	}
	defer store.database.Disconnect(ctx)

	var todo Todo
	err = store.database.Collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&todo)
	if err != nil {
		return Todo{}, err
	}

	todo.Status = !todo.Status
	todo.UpdatedAt = time.Now().UTC()

	_, err = store.database.Collection.ReplaceOne(ctx, bson.D{{"_id", id}}, todo)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (store *TodoStore) Delete(ctx context.Context, id uuid.UUID) error {
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
