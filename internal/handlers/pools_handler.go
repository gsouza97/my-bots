package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/httpserver/service"
	"github.com/gsouza97/my-bots/internal/logger"
)

type PoolsHandler struct {
	poolsService *service.PoolsService
}

func NewPoolsHandler(poolsService *service.PoolsService) *PoolsHandler {
	return &PoolsHandler{
		poolsService: poolsService,
	}
}

func (h *PoolsHandler) GetAllPools(c *gin.Context) {
	pools, err := h.poolsService.GetAll()
	if err != nil {
		logger.Log.Errorf("Error getting pools:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pools"})
		return
	}

	c.JSON(http.StatusOK, pools)
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

// func (h *AlertsHandler) CreateAlert(c *gin.Context) {
// 	var input dto.CreatePriceAlertInput

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	alert, err := h.alertsService.CreateAlert(input)
// 	if err != nil {
// 		logger.Log.Errorf("Error creating alert:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, alert)
// }
