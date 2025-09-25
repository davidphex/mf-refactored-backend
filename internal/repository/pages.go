package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/davidphex/memoryframe-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const COLLECTION_NAME_PAGES = "pages"

type PagesRepository interface {
	GetPages(albumId string) (*[]models.AlbumPage, error)
}

type pagesRepository struct {
	db *mongo.Database
}

func NewPagesRepository(db *mongo.Database) PagesRepository {
	return &pagesRepository{db: db}
}

func (r *pagesRepository) GetPages(albumId string) (*[]models.AlbumPage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	albumObjectId, err := bson.ObjectIDFromHex(albumId)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Fetching pages for album ID: %s\n", albumId)

	var pages []models.AlbumPage
	cursor, err := r.db.Collection(COLLECTION_NAME_PAGES).Find(ctx, bson.M{albumId: albumObjectId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}

	fmt.Printf("Retrieved %d pages for album ID %s\n", len(pages), albumId)

	return &pages, nil
}
