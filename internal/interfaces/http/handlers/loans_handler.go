package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/application/dto"
	"github.com/gsouza97/my-bots/internal/application/services"
	"github.com/gsouza97/my-bots/internal/logger"
)

type LoansHandler struct {
	loansService *services.LoansService
}

func NewLoansHandler(loansService *services.LoansService) *LoansHandler {
	return &LoansHandler{
		loansService: loansService,
	}
}

func (h *LoansHandler) GetAllLoans(c *gin.Context) {
	loans, err := h.loansService.GetAll()
	if err != nil {
		logger.Log.Errorf("Error getting loans:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get loans"})
		return
	}

	c.JSON(http.StatusOK, loans)
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

func (h *LoansHandler) CreateLoan(c *gin.Context) {
	var input dto.CreateLoanInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.loansService.CreateLoan(input)
	if err != nil {
		logger.Log.Errorf("Error creating loan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}
