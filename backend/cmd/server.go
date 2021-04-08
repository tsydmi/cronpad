package cmd

import (
	"context"
	"github.com/ts-dmitry/cronpad/backend/rest"
	"log"
)

func RunApp() error {
	mongoConfig := getMongoConfig()
	client, err := connectToMongo(context.TODO(), mongoConfig)
	if err != nil {
		return err
	}

	defer func() {
		err := client.Disconnect(context.TODO())
		if err == nil {
			log.Println("Connection to MongoDB is closed.")
		}
	}()

	keycloakUrl := getEnv("KEYCLOAK_URL", "http://localhost:8080")
	err = keycloakHealthCheck(keycloakUrl)
	if err != nil {
		return err
	}

	server, err := rest.CreateRestServer(client.Database(mongoConfig.db), keycloakUrl)
	if err != nil {
		return err
	}

	return server.Run()
}
