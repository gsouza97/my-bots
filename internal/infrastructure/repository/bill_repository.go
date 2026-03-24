package repository

import (
	"context"
	"time"

	"github.com/gsouza97/my-bots/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BillRepository interface {
	Save(ctx context.Context, bill *domain.Bill) (*domain.Bill, error)
	FindByID(ctx context.Context, id string) (*domain.Bill, error)
	FindAll(ctx context.Context) ([]*domain.Bill, error)
	FindByMonth(ctx context.Context, month time.Month) ([]*domain.Bill, error)
	FindByPurchaseDateBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*domain.Bill, error)
}

type billRepository struct {
	collection *mongo.Collection
}

func NewBillRepository(db *mongo.Database) BillRepository {
	return &billRepository{
		collection: db.Collection("bills"),
	}
}

func (r *billRepository) Save(ctx context.Context, bill *domain.Bill) (*domain.Bill, error) {
	bill.Timestamp = time.Now()

	if bill.ID.IsZero() {
		bill.ID = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, bill)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (r *billRepository) FindByID(ctx context.Context, id string) (*domain.Bill, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var bill domain.Bill
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&bill)
	if err != nil {
		return nil, err
	}

	return &bill, nil
}

func (r *billRepository) FindAll(ctx context.Context) ([]*domain.Bill, error) {
	var bills []*domain.Bill

	// bson.M{} eh o equivalente a um SELECT * FROM bills
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var bill domain.Bill
		if err := cursor.Decode(&bill); err != nil {
			return nil, err
		}
		bills = append(bills, &bill)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bills, nil
}

func (r *billRepository) FindByPurchaseDateBetween(ctx context.Context, startDate time.Time, endDate time.Time) ([]*domain.Bill, error) {
	var bills []*domain.Bill

	// Select * FROM bills WHERE purchaseDate BETWEEN startDate AND endDate
	cursor, err := r.collection.Find(ctx, bson.M{
		"purchaseDate": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var bill domain.Bill
		if err := cursor.Decode(&bill); err != nil {
			return nil, err
		}
		bills = append(bills, &bill)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bills, nil
}

func (r *billRepository) FindByMonth(ctx context.Context, month time.Month) ([]*domain.Bill, error) {
	year := time.Now().Year()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return r.FindByPurchaseDateBetween(ctx, startOfMonth, endOfMonth)
}
