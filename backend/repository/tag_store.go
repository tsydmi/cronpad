package repository

import (
	"context"
	"github.com/ts-dmitry/cronpad/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type TagStore struct {
	collection   *mongo.Collection
	uuidProvider utils.UuidProvider
}

func CreateTagStore(database *mongo.Database, uuidProvider utils.UuidProvider) *TagStore {
	const collectionName = "tag"

	return &TagStore{collection: database.Collection(collectionName), uuidProvider: uuidProvider}
}

func (t *TagStore) Create(record Tag) (*mongo.InsertOneResult, error) {
	record.ID = t.uuidProvider.New()
	result, err := t.collection.InsertOne(context.TODO(), record)
	return result, err
}

func (t *TagStore) FindAll(userID string) ([]Tag, error) {
	filter := bson.D{{"userid", userID}}
	cursor, err := t.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return getTagResults(cursor)
}

func (t *TagStore) Delete(tagID string, userID string) error {
	filter := bson.D{{"_id", tagID}, {"userid", userID}}
	_, err := t.collection.DeleteOne(context.TODO(), filter)
	return err
}

func getTagResults(cursor *mongo.Cursor) ([]Tag, error) {
	var results = make([]Tag, 0)

	for cursor.Next(context.TODO()) {
		var elem Tag

		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		elem.PrepareToSend()

		results = append(results, elem)
	}

	return results, nil
}
