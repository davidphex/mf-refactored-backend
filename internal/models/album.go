package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Album struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
}
