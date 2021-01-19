package cmd

import (
	"context"
	"github.com/ts-dmitry/cronpad/backend/rest"
	"log"
	"time"
)

const timeout = 2 * time.Minute

func RunApp() {
	mongoConfig := mongoConfig{
		host:     getEnv("MONGO_HOST", "localhost"),
		port:     getEnv("MONGO_PORT", "27017"),
		db:       getEnv("MONGO_DB", "cronpad"),
		username: getEnv("MONGO_USER", "user"),
		password: getEnv("MONGO_PASSWORD", "pwd"),
	}
	keycloakUrl := getEnv("KEYCLOAK_URL", "http://localhost:8080")

	client, err := createMongoClient(mongoConfig)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	authenticator, err := rest.CreateAuthService(keycloakUrl, timeout)
	if err != nil {
		log.Fatal(err)
	}

	server := rest.CreateRestServer(client.Database(mongoConfig.db), authenticator, keycloakUrl)
	server.Run()
}
