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

type AlbumPage struct {
	ID         bson.ObjectID      `bson:"_id" json:"id"`
	AlbumID    bson.ObjectID      `bson:"albumId" json:"albumId" binding:"required"`
	Type       string             `bson:"type" json:"type" binding:"required"`
	PageNumber int                `bson:"pageNumber" json:"pageNumber" binding:"required"`
	Elements   []AlbumPageElement `bson:"elements" json:"elements" binding:"required"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type AlbumPageElement struct {
	ID     bson.ObjectID     `bson:"_id" json:"id"`
	Type   string            `bson:"type" json:"type" binding:"required"`
	Width  float32           `bson:"width" json:"width" binding:"required"`
	Height float32           `bson:"height" json:"height" binding:"required"`
	Top    float32           `bson:"top" json:"top" binding:"required"`
	Left   float32           `bson:"left" json:"left" binding:"required"`
	Style  map[string]string `bson:"style" json:"style" binding:"required"`
	Src    string            `bson:"src,omitempty" json:"src,omitempty"`
	Alt    string            `bson:"alt,omitempty" json:"alt,omitempty"`
}
