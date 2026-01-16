package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartHealthRoutes(r *gin.Engine) {
	r.GET("/health", handleHealthCheck)
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
}
