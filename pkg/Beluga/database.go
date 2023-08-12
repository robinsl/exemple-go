package Beluga

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Database struct {
	database   DatabaseConfiguration
	client     *mongo.Client
	Collection *mongo.Collection
}

func NewDatabase(database DatabaseConfiguration) Database {
	return Database{
		database: database,
	}
}

func (store *Database) Connect(ctx context.Context) error {
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      store.database.Username,
		Password:      store.database.Password,
	}
	clientOpts := options.Client().ApplyURI(store.database.ConnectionString).SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOpts)
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

	store.Collection = client.Database(store.database.DatabaseName).Collection(store.database.DatabaseName)
	log.Print("Connected to database:", store.database.DatabaseName)
	return nil
}

func (store *Database) Disconnect(ctx context.Context) error {
	if err := store.client.Disconnect(ctx); err != nil {
		return err
	}

	log.Print("Disconnected from database:", store.database.DatabaseName)
	return nil
}
