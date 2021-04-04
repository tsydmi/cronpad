package cmd

import (
	"context"
	"fmt"
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

func getMongoConfig() mongoConfig {
	return mongoConfig{
		host:     getEnv("MONGO_HOST", "localhost"),
		port:     getEnv("MONGO_PORT", "27017"),
		db:       getEnv("MONGO_DB", "cronpad"),
		username: getEnv("MONGO_USER", "user"),
		password: getEnv("MONGO_PASSWORD", "pwd"),
	}
}

func connectToMongo(ctx context.Context, config mongoConfig) (*mongo.Client, error) {
	credential := options.Credential{
		AuthSource: config.db,
		Username:   config.username,
		Password:   config.password,
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + config.host + ":" + config.port).SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connection to MongoDB is opened.")

	return client, err
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
