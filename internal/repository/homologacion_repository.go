package repository

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HomologacionRepository interface {
	FindOne(ctx context.Context) (*domain.HomologacionConfigParams, error)
	FindAll(ctx context.Context) ([]*domain.HomologacionConfigParams, error)
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

func (r *homologacionRepositoryRepository) FindAll(ctx context.Context) ([]*domain.HomologacionConfigParams, error) {
	var homologacionConfigParams []*domain.HomologacionConfigParams

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var params domain.HomologacionConfigParams
		if err := cursor.Decode(&params); err != nil {
			return nil, err
		}
		homologacionConfigParams = append(homologacionConfigParams, &params)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return homologacionConfigParams, nil
}
