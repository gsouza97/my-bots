package repository

import (
	"context"

	"github.com/gsouza97/my-bots/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoanRepository interface {
	FindById(ctx context.Context, id string) (*domain.Loan, error)
	FindAll(ctx context.Context) ([]*domain.Loan, error)
	Update(ctx context.Context, loan *domain.Loan) error
	Create(ctx context.Context, loan *domain.Loan) (*domain.Loan, error)
}

type loanRepository struct {
	collection *mongo.Collection
}

func NewLoanRepository(db *mongo.Database) LoanRepository {
	return &loanRepository{
		collection: db.Collection("loans"),
	}
}

func (r *loanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": loan.ID}, bson.M{"$set": loan})
	return err
}

func (r *loanRepository) Create(ctx context.Context, loan *domain.Loan) (*domain.Loan, error) {
	if loan.ID.IsZero() {
		loan.ID = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (r *loanRepository) FindAll(ctx context.Context) ([]*domain.Loan, error) {
	var loans []*domain.Loan
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var loan domain.Loan
		if err := cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	return loans, nil
}

func (r *loanRepository) FindById(ctx context.Context, id string) (*domain.Loan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var loan domain.Loan
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&loan)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}
