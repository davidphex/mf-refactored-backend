package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Photo struct {
	ID         bson.ObjectID `bson:"_id" json:"id"`
	AlbumID    bson.ObjectID `bson:"albumId" json:"albumId"`
	Source     string        `bson:"source" json:"source"`
	Name       string        `bson:"name" json:"name"`
	Size       int64         `bson:"size" json:"size"` // Size in bytes
	Resolution string        `bson:"resolution" json:"resolution"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time     `bson:"updatedAt" json:"updatedAt"`
}
