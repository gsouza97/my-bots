package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/handlers"
)

type LoansRoutes struct {
	loansHandler *handlers.LoansHandler
}

func NewLoansRoutes(loansHandler *handlers.LoansHandler) *LoansRoutes {
	return &LoansRoutes{
		loansHandler: loansHandler,
	}
}

func (rt *LoansRoutes) StartLoansRoutes(r *gin.Engine, middleware ...gin.HandlerFunc) {
	r.GET("/loans", append(middleware, rt.loansHandler.GetAllLoans)...)
	// r.PATCH("/loans/:id", rt.loansHandler.UpdateLoan)
	// r.POST("/loans", rt.loansHandler.CreateLoan)
}
