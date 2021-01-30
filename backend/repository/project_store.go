package repository

import (
	"context"
	"github.com/ts-dmitry/cronpad/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectStore struct {
	collection   *mongo.Collection
	uuidProvider utils.UuidProvider
}

func CreateProjectStore(database *mongo.Database, uuidProvider utils.UuidProvider) *ProjectStore {
	const collectionName = "project"

	return &ProjectStore{collection: database.Collection(collectionName), uuidProvider: uuidProvider}
}

func (t *ProjectStore) Create(record Project) (*mongo.InsertOneResult, error) {
	record.ID = t.uuidProvider.New()
	result, err := t.collection.InsertOne(context.TODO(), record)

	return result, err
}

func (t *ProjectStore) GetProjectWithUsersByID(projectID string) (Project, error) {
	filter := bson.D{{"_id", projectID}}
	var project Project
	err := t.collection.FindOne(context.TODO(), filter).Decode(&project)

	return project, err
}

func (t *ProjectStore) FindAllProjectsByUser(userID string) ([]Project, error) {
	filter := bson.M{"users": userID}
	cursor, err := t.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return getProjectResults(cursor)
}

func (t *ProjectStore) FindAllActiveProjectsByUser(userID string) (Projects, error) {
	filter := bson.D{bson.E{Key: "users", Value: userID}, isActive}
	cursor, err := t.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return getProjectResults(cursor)
}

func (t *ProjectStore) Search(form ProjectSearchForm) ([]Project, error) {
	filters := bson.D{isActive}

	if len(form.Name) > 0 {
		filters = append(filters, bson.E{Key: "name", Value: likeRegex(form.Name)})
	}

	if len(form.Description) > 0 {
		filters = append(filters, bson.E{Key: "description", Value: likeRegex(form.Description)})
	}

	if form.Users != nil && len(form.Users) > 0 {
		filters = append(filters, bson.E{Key: "users", Value: bson.D{
			{"$all", form.Users},
		}})
	}

	cursor, err := t.collection.Find(context.TODO(), filters)
	if err != nil {
		return nil, err
	}

	return getProjectResults(cursor)
}

func (t *ProjectStore) Update(project Project) (string, error) {
	filter := bson.D{{"_id", project.ID}}
	var updatedProject Project
	err := t.collection.FindOneAndReplace(context.TODO(), filter, project).Decode(&updatedProject)

	return updatedProject.ID, err
}

func (t *ProjectStore) Delete(projectID string) error {
	filter := bson.D{{"_id", projectID}}
	update := bson.D{{"$set", bson.D{{"active", false}}}}

	_, err := t.collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func getProjectResults(cursor *mongo.Cursor) ([]Project, error) {
	var results = make([]Project, 0)

	for cursor.Next(context.TODO()) {
		var elem Project

		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	return results, nil
}

type ProjectSearchForm struct {
	Name        string
	Description string
	Users       []string
}
