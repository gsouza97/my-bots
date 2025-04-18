package repository

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PriceAlertRepository interface {
	FindAllByActiveIsTrue(ctx context.Context) ([]*domain.PriceAlert, error)
	Update(ctx context.Context, alert *domain.PriceAlert) error
}

type priceAlertRepository struct {
	collection *mongo.Collection
}

func NewPriceAlertRepository(db *mongo.Database) PriceAlertRepository {
	return &priceAlertRepository{
		collection: db.Collection("alerts"),
	}
}

func (r *priceAlertRepository) FindAllByActiveIsTrue(ctx context.Context) ([]*domain.PriceAlert, error) {
	var alerts []*domain.PriceAlert

	cursor, err := r.collection.Find(ctx, bson.M{"active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var alert domain.PriceAlert
		if err := cursor.Decode(&alert); err != nil {
			return nil, err
		}
		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

func (r *priceAlertRepository) Update(ctx context.Context, alert *domain.PriceAlert) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": alert.ID}, bson.M{"$set": alert})
	return err
}
