package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartHealthRoutes(r *gin.Engine) {
	r.GET("/healthz", handleHealthCheck)
	r.HEAD("/health", handleUptimeHealthCheck)
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
}

func handleUptimeHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
}
