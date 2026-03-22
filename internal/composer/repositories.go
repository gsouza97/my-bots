package composer

import (
	"github.com/gsouza97/my-bots/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoriesComposer struct {
	BillRepository         repository.BillRepository
	PriceAlertRepository   repository.PriceAlertRepository
	PoolRepository         repository.PoolRepository
	HomologacionRepository repository.HomologacionRepository
	UserRepository         repository.UserRepository
	LoanRepository         repository.LoanRepository
}

func NewRepositoriesComposer(db *mongo.Database) *RepositoriesComposer {
	return &RepositoriesComposer{
		BillRepository:         repository.NewBillRepository(db),
		PriceAlertRepository:   repository.NewPriceAlertRepository(db),
		PoolRepository:         repository.NewPoolRepository(db),
		HomologacionRepository: repository.NewHomologacionRepository(db),
		UserRepository:         repository.NewUserRepository(db),
		LoanRepository:         repository.NewLoanRepository(db),
	}
}
