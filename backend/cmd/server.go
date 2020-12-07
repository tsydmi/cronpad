package cmd

import (
	"context"
	"github.com/ts-dmitry/cronpad/backend/rest"
	"log"
	"time"
)

func RunApp() {
	client, mongoConfig, err := createMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	server := rest.CreateRestServer(client.Database(mongoConfig.db), "http://localhost:8080")
	server.Run()
}
