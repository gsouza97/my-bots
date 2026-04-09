package services

import (
	"context"
	"time"

	"github.com/gsouza97/my-bots/internal/application/dto"
	"github.com/gsouza97/my-bots/internal/domain"
	"github.com/gsouza97/my-bots/internal/infrastructure/repository"
	"github.com/gsouza97/my-bots/internal/logger"
)

type BillsService struct {
	billRepository repository.BillRepository
}

func NewBillsService(billRepository repository.BillRepository) *BillsService {
	return &BillsService{
		billRepository: billRepository,
	}
}

func (s *BillsService) GetAll() ([]*domain.Bill, error) {
	ctx := context.Background()

	bills, err := s.billRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return bills, nil
}

func (s *BillsService) UpdateBill(id string, input dto.UpdateBillInput) (*domain.Bill, error) {
	ctx := context.Background()
	bill, err := s.billRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bill.Description = input.Description
	bill.Amount = input.Amount
	bill.Category = input.Category
	bill.PurchaseDate, err = time.Parse("2006-01-02", input.PurchaseDate)
	if err != nil {
		return nil, err
	}

	err = s.billRepository.Update(ctx, bill)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (uc *BillsService) CreateBill(input dto.CreateBillInput) (*domain.Bill, error) {
	ctx := context.Background()

	parsedPurchaseDate, err := time.Parse("2006-01-02", input.PurchaseDate)
	logger.Log.Infof("Parsed purchase date: %v", parsedPurchaseDate)
	if err != nil {
		return nil, err
	}

	bill := &domain.Bill{
		Description:  input.Description,
		Amount:       input.Amount,
		Category:     input.Category,
		PurchaseDate: parsedPurchaseDate,
	}

	createdBill, err := uc.billRepository.Save(ctx, bill)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Bill created: %+v", createdBill)

	return createdBill, nil
}

func (s *BillsService) DeleteBill(id string) error {
	ctx := context.Background()
	return s.billRepository.Delete(ctx, id)
}
