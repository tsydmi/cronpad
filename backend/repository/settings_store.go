package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SettingsStore struct {
	collection *mongo.Collection
}

func CreateSettingsStore(database *mongo.Database) *SettingsStore {
	const collectionName = "settings"

	return &SettingsStore{collection: database.Collection(collectionName)}
}

func (s *SettingsStore) FindByUser(userID string) (*Settings, error) {
	var settings *Settings
	filter := bson.D{{Key: "user", Value: userID}}
	err := s.collection.FindOne(context.TODO(), filter).Decode(&settings)
	return settings, err
}

func (s *SettingsStore) CreateOrUpdate(settings Settings) error {
	filter := bson.D{{Key: "user", Value: settings.UserID}}
	opts := options.Replace().SetUpsert(true)

	_, err := s.collection.ReplaceOne(context.TODO(), filter, settings, opts)
	return err
}
