package todo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"goexemples/pkg/Beluga"
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
