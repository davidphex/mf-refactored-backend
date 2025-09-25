package repository

import (
	"context"
	"time"

	"github.com/davidphex/memoryframe-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AlbumRepository interface {
	GetAllAlbums() (*[]models.Album, error)
	GetAlbumByID(id string) (*models.Album, error)
}

type repository struct {
	db *mongo.Database
}

func NewAlbumRepository(db *mongo.Database) AlbumRepository {
	return &repository{db: db}
}

func (r *repository) GetAllAlbums() (*[]models.Album, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.db.Collection("Albums").Find(ctx, bson.D{})
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

func (r *repository) GetAlbumByID(id string) (*models.Album, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bsonID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var album models.Album
	err = r.db.Collection("Albums").FindOne(ctx, bson.M{"_id": bsonID}).Decode(&album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}
