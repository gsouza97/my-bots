package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/dto"
	"github.com/gsouza97/my-bots/internal/logger"
	"github.com/gsouza97/my-bots/internal/usecase"
)

type AlertsHandler struct {
	alertsUseCase *usecase.AlertsUseCase
}

func NewAlertsHandler(alertsUseCase *usecase.AlertsUseCase) *AlertsHandler {
	return &AlertsHandler{
		alertsUseCase: alertsUseCase,
	}
}

func (h *AlertsHandler) GetAllAlerts(c *gin.Context) {
	alerts, err := h.alertsUseCase.GetAll()
	if err != nil {
		logger.Log.Errorf("Error getting alerts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alerts"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (h *AlertsHandler) UpdateAlert(c *gin.Context) {
	alertID := c.Param("id")
	var input dto.UpdatePriceAlertInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.alertsUseCase.UpdateAlert(alertID, input)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alert)
}
