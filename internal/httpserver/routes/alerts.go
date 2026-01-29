package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/handlers"
)

type AlertsRoutes struct {
	alertsHandler *handlers.AlertsHandler
}

func NewAlertsRoutes(alertsHandler *handlers.AlertsHandler) *AlertsRoutes {
	return &AlertsRoutes{
		alertsHandler: alertsHandler,
	}
}

func (rt *AlertsRoutes) StartAlertsRoutes(r *gin.Engine) {
	r.GET("/alerts", rt.alertsHandler.GetAllAlerts)
	r.PATCH("/alerts/:id", rt.alertsHandler.UpdateAlert)
	r.POST("/alerts", rt.alertsHandler.CreateAlert)
}
