package main

import (
	"context"
	"goexemples/internal/exemple/api"
	"goexemples/internal/exemple/config"
	"goexemples/internal/exemple/store"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	store := store.NewBookStore(cfg.Database)
	server := api.NewServer(cfg.HttpServer, store)
	server.Serve(ctx)
}
