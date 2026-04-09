package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/application/dto"
	"github.com/gsouza97/my-bots/internal/application/services"
	"github.com/gsouza97/my-bots/internal/logger"
)

type BillsHandler struct {
	billsService *services.BillsService
}

func NewBillsHandler(billsService *services.BillsService) *BillsHandler {
	return &BillsHandler{
		billsService: billsService,
	}
}

func (h *BillsHandler) GetAllBills(c *gin.Context) {
	bills, err := h.billsService.GetAll()
	if err != nil {
		logger.Log.Errorf("Error getting bills:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bills"})
		return
	}

	c.JSON(http.StatusOK, bills)
}

// func (h *AlertsHandler) UpdateAlert(c *gin.Context) {
// 	alertID := c.Param("id")
// 	var input dto.UpdatePriceAlertInput

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	alert, err := h.alertsService.UpdateAlert(alertID, input)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, alert)
// }

func (h *BillsHandler) CreateBill(c *gin.Context) {
	var input dto.CreateBillInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bill, err := h.billsService.CreateBill(input)
	if err != nil {
		logger.Log.Errorf("Error creating bill:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bill)
}

// func (h *AlertsHandler) DeleteAlert(c *gin.Context) {
// 	alertID := c.Param("id")

// 	err := h.alertsService.DeleteAlert(alertID)
// 	if err != nil {
// 		logger.Log.Errorf("Error deleting alert:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})
// }
