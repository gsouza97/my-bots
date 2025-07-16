package repository

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HomologacionRepository interface {
	FindOne(ctx context.Context) (*domain.HomologacionConfigParams, error)
}

type homologacionRepositoryRepository struct {
	collection *mongo.Collection
}

func NewHomologacionRepository(db *mongo.Database) HomologacionRepository {
	return &homologacionRepositoryRepository{
		collection: db.Collection("homologacion"),
	}
}

func (r *homologacionRepositoryRepository) FindOne(ctx context.Context) (*domain.HomologacionConfigParams, error) {
	var homologacionConfigParams *domain.HomologacionConfigParams

	err := r.collection.FindOne(ctx, bson.M{}).Decode(&homologacionConfigParams)
	if err != nil {
		return nil, err
	}

	return homologacionConfigParams, nil
}
