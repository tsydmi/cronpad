package cmd

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type mongoConfig struct {
	host     string
	port     string
	db       string
	username string
	password string
}

func createMongoClient(config mongoConfig) (*mongo.Client, error) {
	credential := options.Credential{
		AuthSource: config.db,
		Username:   config.username,
		Password:   config.password,
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + config.host + ":" + config.port).SetAuth(credential)

	client, err := mongo.NewClient(clientOptions)
	return client, err
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
