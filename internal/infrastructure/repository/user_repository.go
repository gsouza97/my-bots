package repository

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindAdminUser(ctx context.Context) (*domain.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) FindAdminUser(ctx context.Context) (*domain.User, error) {

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
