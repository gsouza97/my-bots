package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartHealthRoutes(r *gin.Engine) {
	r.GET("/healthz", handleHealthCheck)
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
}
