package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
