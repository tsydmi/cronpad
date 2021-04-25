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

func (t *DayStore) Search(form DaySearchForm) ([]Day, error) {
	filters := bson.D{}
	if !form.From.IsZero() || !form.To.IsZero() {
		dateFilters := bson.D{}
		if !form.From.IsZero() {
			dateFilters = append(dateFilters, bson.E{Key: "$gte", Value: form.From})
		}
		if !form.To.IsZero() {
			dateFilters = append(dateFilters, bson.E{Key: "$lte", Value: form.To})
		}

		filters = append(filters, bson.E{Key: "date", Value: dateFilters})
	}

	if len(form.UserID) > 0 {
		filters = append(filters, bson.E{Key: "userid", Value: form.UserID})
	}

	if len(form.UserIDs) > 0 {
		filters = append(filters, bson.E{Key: "userid", Value: bson.M{"$in": form.UserIDs}})
	}

	cursor, err := t.collection.Find(context.TODO(), filters)
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

func (t *DayStore) GetUsedNames(userID string, tagID string, from time.Time, to time.Time) ([]string, error) {
	match := []bson.M{
		{"$match": bson.M{
			"date": bson.M{
				"$gte": from,
				"$lte": to,
			},
			"userid": userID,
			"events": bson.M{
				"$elemMatch": bson.M{"tagid": tagID},
			},
		}},
		{"$sort": bson.M{"date": -1}},
		{"$unwind": bson.M{"path": "$events"}},
		{"$group": bson.M{
			"_id": bson.M{
				"name":  "$events.name",
				"tagid": tagID,
			},
			"date": bson.M{"$first": "$date"},
		}},
		{"$project": bson.M{"name": "$_id.name"}},
		{"$limit": 5},
	}

	cursor, err := t.collection.Aggregate(context.TODO(), match)
	if err != nil {
		return nil, err
	}

	return getUsedEventNamesResults(cursor)
}

func getDayResults(cursor *mongo.Cursor) ([]Day, error) {
	var results = make([]Day, 0)

	for cursor.Next(context.TODO()) {
		var elem Day

		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	return results, nil
}

type DaySearchForm struct {
	UserID  string
	From    time.Time
	To      time.Time
	UserIDs []string
}

type usedEventNames struct {
	Name string
}

func getUsedEventNamesResults(cursor *mongo.Cursor) ([]string, error) {
	var results = make([]string, 0)

	for cursor.Next(context.TODO()) {
		var elem usedEventNames

		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem.Name)
	}

	return results, nil
}
