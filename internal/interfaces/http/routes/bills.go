package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/interfaces/http/handlers"
)

type BillsRoutes struct {
	billsHandler *handlers.BillsHandler
}

func NewBillsRoutes(billsHandler *handlers.BillsHandler) *BillsRoutes {
	return &BillsRoutes{
		billsHandler: billsHandler,
	}
}

func (rt *BillsRoutes) StartBillsRoutes(r *gin.Engine, middleware ...gin.HandlerFunc) {
	r.GET("/bills", append(middleware, rt.billsHandler.GetAllBills)...)
}
