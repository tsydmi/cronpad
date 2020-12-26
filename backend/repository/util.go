package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

func likeRegex(s string) primitive.Regex {
	return primitive.Regex{Pattern: ".*" + s + ".*", Options: ""}
}
