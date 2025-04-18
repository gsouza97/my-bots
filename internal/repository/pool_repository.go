package repository

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PoolRepository interface {
	Save(ctx context.Context, pool *domain.Pool) (*domain.Pool, error)
	FindByID(ctx context.Context, id string) (*domain.Pool, error)
	FindAll(ctx context.Context) ([]*domain.Pool, error)
	FindAllByActiveIsTrue(ctx context.Context) ([]*domain.Pool, error)
	Update(ctx context.Context, pool *domain.Pool) error
}

type poolRepository struct {
	collection *mongo.Collection
}

func NewPoolRepository(db *mongo.Database) PoolRepository {
	return &poolRepository{
		collection: db.Collection("pools"),
	}
}

func (r *poolRepository) Save(ctx context.Context, pool *domain.Pool) (*domain.Pool, error) {

	if pool.ID.IsZero() {
		pool.ID = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (r *poolRepository) FindByID(ctx context.Context, id string) (*domain.Pool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pool domain.Pool
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&pool)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func (r *poolRepository) FindAll(ctx context.Context) ([]*domain.Pool, error) {
	var pools []*domain.Pool

	// bson.M{} eh o equivalente a um SELECT * FROM pools
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var pool domain.Pool
		if err := cursor.Decode(&pool); err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return pools, nil
}

func (r *poolRepository) FindAllByActiveIsTrue(ctx context.Context) ([]*domain.Pool, error) {
	var pools []*domain.Pool

	cursor, err := r.collection.Find(ctx, bson.M{"active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var pool domain.Pool
		if err := cursor.Decode(&pool); err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	return pools, nil
}

func (r *poolRepository) Update(ctx context.Context, pool *domain.Pool) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": pool.ID}, bson.M{"$set": pool})
	return err
}
