package cmd

import (
	"context"
	"fmt"
	"github.com/ts-dmitry/cronpad/backend/rest"
	"time"
)

const keycloakTimeout = 2 * time.Minute

func RunApp() error {
	mongoConfig := getMongoConfig()
	client, err := connectToMongo(context.TODO(), mongoConfig)
	if err != nil {
		return err
	}

	defer func() {
		err := client.Disconnect(context.TODO())
		if err == nil {
			fmt.Println("Connection to MongoDB is closed.")
		}
	}()

	keycloakUrl := getEnv("KEYCLOAK_URL", "http://localhost:8080")
	authenticator, err := rest.CreateJwtAuthService(keycloakUrl, keycloakTimeout)
	if err != nil {
		return err
	}

	server := rest.CreateRestServer(client.Database(mongoConfig.db), authenticator, keycloakUrl)
	return server.Run()
}
