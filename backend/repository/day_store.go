package repository

import (
	"context"
	"github.com/ts-dmitry/cronpad/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type DayStore struct {
	collection   *mongo.Collection
	uuidProvider utils.UuidProvider
}

func CreateDayStore(database *mongo.Database, uuidProvider utils.UuidProvider) *DayStore {
	const collectionName = "day"

	return &DayStore{collection: database.Collection(collectionName), uuidProvider: uuidProvider}
}

func (t *DayStore) Create(day Day) (*mongo.InsertOneResult, error) {
	day.ID = t.uuidProvider.New()
	result, err := t.collection.InsertOne(context.TODO(), day)
	return result, err
}

func (t *DayStore) FindByEventID(eventID string, userID string) (Day, error) {
	filter := bson.M{"events": bson.M{"$elemMatch": bson.M{"_id": eventID}}, "userid": userID}
	var day Day
	err := t.collection.FindOne(context.TODO(), filter).Decode(&day)
	return day, err
}

func (t *DayStore) FindByDate(date time.Time, userID string) (Day, error) {
	dateWithoutTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	filter := bson.M{"date": dateWithoutTime, "userid": userID}
	var day Day
	err := t.collection.FindOne(context.TODO(), filter).Decode(&day)
	return day, err
}

func (t *DayStore) FindByDateRange(from time.Time, to time.Time, userID string) ([]Day, error) {
	filter := bson.M{"date": bson.M{"$gte": from, "$lte": to}, "userid": userID}
	cursor, err := t.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return getDayResults(cursor)
}

func (t *DayStore) Update(day Day) (Day, error) {
	filter := bson.D{{"_id", day.ID}, {"userid", day.UserID}}
	var replacedRecord Day
	err := t.collection.FindOneAndReplace(context.TODO(), filter, day).Decode(&replacedRecord)
	return day, err
}

func getDayResults(cursor *mongo.Cursor) ([]Day, error) {
	var results = make([]Day, 0)

	for cursor.Next(context.TODO()) {
		var elem Day

		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		elem.PrepareToSend()

		results = append(results, elem)
	}

	return results, nil
}
