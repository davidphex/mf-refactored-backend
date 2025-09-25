package repository

import (
	"context"
	"time"

	"github.com/davidphex/memoryframe-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const COLLECTION_NAME = "albums"

type AlbumRepository interface {
	GetAll() (*[]models.Album, error)
	GetById(id string) (*models.Album, error)
	Insert(album *models.Album) error
	Update(album *models.Album) error
	Delete(id string) error
}

type repository struct {
	db *mongo.Database
}

func NewAlbumRepository(db *mongo.Database) AlbumRepository {
	return &repository{db: db}
}

func (r *repository) GetAll() (*[]models.Album, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.db.Collection(COLLECTION_NAME).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var albums []models.Album
	if err = cursor.All(ctx, &albums); err != nil {
		return nil, err
	}

	return &albums, nil
}

func (r *repository) GetById(id string) (*models.Album, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bsonID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var album models.Album
	err = r.db.Collection(COLLECTION_NAME).FindOne(ctx, bson.M{"_id": bsonID}).Decode(&album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func (r *repository) Insert(album *models.Album) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set ID and timestamps
	album.ID = bson.NewObjectID()
	album.CreatedAt = time.Now()
	album.UpdatedAt = time.Now()

	_, err := r.db.Collection(COLLECTION_NAME).InsertOne(ctx, album)
	return err
}

func (r *repository) Update(album *models.Album) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the UpdatedAt timestamp
	album.UpdatedAt = time.Now()

	_, err := r.db.Collection(COLLECTION_NAME).UpdateOne(ctx, bson.M{"_id": album.ID}, bson.M{"$set": album})
	return err
}

func (r *repository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bsonID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(COLLECTION_NAME).DeleteOne(ctx, bson.M{"_id": bsonID})
	return err
}
