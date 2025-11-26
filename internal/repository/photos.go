package repository

import (
	"context"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/davidphex/memoryframe-backend/internal/models"
)

const COLLECTION_NAME_PHOTOS = "photos"

type PhotoRepository interface {
	// Basic CRUD operations
	Get(photoId string) (*models.Photo, error)
	GetByAlbumId(albumId string) ([]*models.Photo, error)
	Insert(file multipart.File, fileName string, albumId string) (string, error) // returns the photo ID
	Update(photo *models.Photo) error
	Delete(photoId string) error
}

type photoRepository struct {
	db  *mongo.Database
	cld *cloudinary.Cloudinary
}

func NewPhotoRepository(db *mongo.Database, cld *cloudinary.Cloudinary) PhotoRepository {
	return &photoRepository{db: db, cld: cld}
}

func (r *photoRepository) Get(photoId string) (*models.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	photoObjectId, err := bson.ObjectIDFromHex(photoId)
	if err != nil {
		return nil, err
	}

	var photo models.Photo
	err = r.db.Collection(COLLECTION_NAME_PHOTOS).FindOne(ctx, bson.M{"_id": photoObjectId}).Decode(&photo)
	if err != nil {
		return nil, err
	}

	return &photo, nil
}

func (r *photoRepository) GetByAlbumId(albumId string) ([]*models.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	albumObjectId, err := bson.ObjectIDFromHex(albumId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"albumId": albumObjectId}

	cursor, err := r.db.Collection(COLLECTION_NAME_PHOTOS).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Transform the cursor results into a slice of Photo models
	var photos []*models.Photo
	if err = cursor.All(ctx, &photos); err != nil {
		return nil, err
	}

	return photos, nil
}

func (r *photoRepository) Insert(file multipart.File, fileName string, albumId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	photoObjectId := bson.NewObjectID()
	photoId := photoObjectId.Hex()

	albumObjectId, err := bson.ObjectIDFromHex(albumId)
	if err != nil {
		return "", err
	}

	// Upload the photo to Cloudinary
	resp, err := r.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: photoId,
		Folder:   albumId,
	})
	if err != nil {
		return "", err
	}

	photo := &models.Photo{
		ID:         photoObjectId,
		AlbumID:    albumObjectId,
		Source:     resp.SecureURL,
		Name:       fileName,
		Resolution: strconv.Itoa(resp.Width) + "x" + strconv.Itoa(resp.Height),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = r.db.Collection(COLLECTION_NAME_PHOTOS).InsertOne(ctx, photo)
	if err != nil {
		return "", err
	}

	return photoId, nil
}

func (r *photoRepository) Update(photo *models.Photo) error {
	// TODO: Implement photo update logic
	return nil
}

func (r *photoRepository) Delete(photoId string) error {
	// TODO: Implement photo deletion logic
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	photoObjectId, err := bson.ObjectIDFromHex(photoId)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(COLLECTION_NAME_PHOTOS).DeleteOne(ctx, bson.M{"_id": photoObjectId})
	if err != nil {
		return err
	}

	// Delete the photo from Cloudinary
	_, err = r.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: photoId,
	})
	if err != nil {
		return err
	}

	return nil
}
