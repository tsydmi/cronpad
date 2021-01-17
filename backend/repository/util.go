package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var isActive = bson.E{Key: "active", Value: bson.D{{Key: "$not", Value: bson.D{{Key: "$eq", Value: false}}}}}

func likeRegex(s string) primitive.Regex {
	return primitive.Regex{Pattern: ".*" + s + ".*", Options: ""}
}
