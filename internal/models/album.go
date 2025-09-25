package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Album struct {
	ID          bson.ObjectID   `bson:"_id" json:"id"`
	Title       string          `bson:"title" json:"title" binding:"required"`
	Description string          `bson:"description" json:"description" binding:"required"`
	Thumbnail   string          `bson:"thumbnail" json:"thumbnail" binding:"required,url"`
	CreatorID   bson.ObjectID   `bson:"creatorId" json:"creatorId" binding:"required"`
	AdminsID    []bson.ObjectID `bson:"adminsId" json:"adminsId" binding:"required"`
	MembersID   []bson.ObjectID `bson:"membersId" json:"membersId" binding:"required"`
	PhotosID    []bson.ObjectID `bson:"photos" json:"photos" binding:"required"`
	PagesID     []bson.ObjectID `bson:"pages" json:"pages" binding:"required"`
	Type        string          `bson:"type" json:"type" binding:"required"`
	ChildAlbums []Album         `bson:"childAlbums" json:"childAlbums"`
	CreatedAt   time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time       `bson:"updatedAt" json:"updatedAt"`
}
