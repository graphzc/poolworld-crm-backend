package repository

import (
	"context"
	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	userModel := &model.User{
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, userModel)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(bson.ObjectID).Hex()
	return nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var userModel model.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&userModel)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        userModel.ID.Hex(),
		Username:  userModel.Username,
		Password:  userModel.Password,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}, nil
}
