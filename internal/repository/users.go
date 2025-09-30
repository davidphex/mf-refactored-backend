package repository

import (
	"context"
	"time"

	"github.com/davidphex/memoryframe-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const COLLECTION_NAME_USERS = "users"

type UserRepository interface {
	GetById(id string) (*models.User, error)
	GetByMail(email string) (*models.User, error)
	Insert(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetById(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userObjectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.Collection(COLLECTION_NAME_USERS).FindOne(ctx, bson.M{"_id": userObjectId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByMail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var user models.User
	err := r.db.Collection(COLLECTION_NAME_USERS).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Insert(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Check for existing user with the same email
	existingUser, err := r.GetByMail(user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if existingUser != nil {
		return mongo.ErrNilDocument
	}

	// Generate a new ObjectID for the user
	user.ID = bson.NewObjectID()

	_, err = r.db.Collection(COLLECTION_NAME_USERS).InsertOne(ctx, user)
	return err
}

func (r *userRepository) Update(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := r.db.Collection(COLLECTION_NAME_USERS).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (r *userRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userObjectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(COLLECTION_NAME_USERS).DeleteOne(ctx, bson.M{"_id": userObjectId})
	return err
}
