package repository

import (
	"context"
	"time"

	"github.com/davidphex/memoryframe-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const COLLECTION_NAME_PAGES = "pages"

type PagesRepository interface {
	// Basic CRUD operations
	Get(pageId string) (*models.AlbumPage, error)
	Insert(page *models.AlbumPage) error
	Update(page *models.AlbumPage) error
	Delete(pageId string) error

	// Custom query
	GetByAlbumId(albumId string) (*[]models.AlbumPage, error)
}

type pagesRepository struct {
	db *mongo.Database
}

func NewPagesRepository(db *mongo.Database) PagesRepository {
	return &pagesRepository{db: db}
}

func (r *pagesRepository) GetByAlbumId(albumId string) (*[]models.AlbumPage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	albumObjectId, err := bson.ObjectIDFromHex(albumId)
	if err != nil {
		return nil, err
	}

	var pages []models.AlbumPage
	cursor, err := r.db.Collection(COLLECTION_NAME_PAGES).Find(ctx, bson.M{"albumId": albumObjectId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}

	return &pages, nil
}

func (r *pagesRepository) Get(pageId string) (*models.AlbumPage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pageObjectId, err := bson.ObjectIDFromHex(pageId)
	if err != nil {
		return nil, err
	}

	var page models.AlbumPage
	err = r.db.Collection(COLLECTION_NAME_PAGES).FindOne(ctx, bson.M{"_id": pageObjectId}).Decode(&page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

func (r *pagesRepository) Insert(page *models.AlbumPage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page.ID = bson.NewObjectID()
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()

	_, err := r.db.Collection(COLLECTION_NAME_PAGES).InsertOne(ctx, page)
	return err
}

func (r *pagesRepository) Update(page *models.AlbumPage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page.UpdatedAt = time.Now()

	_, err := r.db.Collection(COLLECTION_NAME_PAGES).UpdateOne(ctx, bson.M{"_id": page.ID}, bson.M{"$set": page})
	return err
}

func (r *pagesRepository) Delete(pageId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pageObjectId, err := bson.ObjectIDFromHex(pageId)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(COLLECTION_NAME_PAGES).DeleteOne(ctx, bson.M{"_id": pageObjectId})
	return err
}
