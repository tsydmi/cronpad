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

func (t *TagStore) Create(tag Tag) (*mongo.InsertOneResult, error) {
	tag.ID = t.uuidProvider.New()
	result, err := t.collection.InsertOne(context.TODO(), tag)
	return result, err
}

func (t *TagStore) FindAll() ([]Tag, error) {
	return t.findAll(bson.D{})
}

func (t *TagStore) FindAllActive() ([]Tag, error) {
	return t.findAll(bson.D{isActive})
}

func (t *TagStore) FindAllBaseAndProjectActiveTags(projectIDs []string) ([]Tag, error) {
	filter := bson.D{isActive,
		{Key: "$or", Value: []bson.M{
			{"basic": true},
			{"project": bson.M{"$in": projectIDs}},
		}},
	}

	return t.findAll(filter)
}

func (t *TagStore) findAll(filter bson.D) ([]Tag, error) {
	cursor, err := t.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return getTagResults(cursor)
}

func (t *TagStore) Update(tag Tag) (string, error) {
	filter := bson.D{
		{Key: "_id", Value: tag.ID},
		{Key: "basic", Value: false},
	}
	return t.update(tag, filter)
}

func (t *TagStore) UpdateBaseTag(tag Tag) (string, error) {
	filter := bson.D{
		{Key: "_id", Value: tag.ID},
		{Key: "basic", Value: true},
	}
	return t.update(tag, filter)
}

func (t *TagStore) update(tag Tag, filter bson.D) (string, error) {
	var updatedTag Tag
	err := t.collection.FindOneAndReplace(context.TODO(), filter, tag).Decode(&updatedTag)

	return updatedTag.ID, err
}
func (t *TagStore) Delete(tagID string) error {
	filter := bson.D{
		{"_id", tagID},
		{Key: "basic", Value: false},
	}
	return t.delete(filter)
}

func (t *TagStore) DeleteByProjectID(tagID string, projectIDs []string) error {
	filter := bson.D{
		{Key: "_id", Value: tagID},
		{Key: "project", Value: bson.M{"$in": projectIDs}},
		{Key: "basic", Value: false},
	}
	return t.delete(filter)
}

func (t *TagStore) DeleteBaseTag(tagID string) error {
	filter := bson.D{
		{Key: "_id", Value: tagID},
		{Key: "basic", Value: true},
	}
	return t.delete(filter)
}

func (t *TagStore) delete(filter bson.D) error {
	update := bson.D{{"$set", bson.D{{"active", false}}}}

	_, err := t.collection.UpdateOne(context.TODO(), filter, update)

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
