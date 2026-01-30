package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/dto"
	"github.com/gsouza97/my-bots/internal/httpserver/service"
	"github.com/gsouza97/my-bots/internal/logger"
)

type LoginHandler struct {
	authService *service.AuthService
}

func NewLoginHandler(authService *service.AuthService) *LoginHandler {
	return &LoginHandler{
		authService: authService,
	}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var input dto.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Authenticate(input)
	if err != nil {
		logger.Log.Errorf("Error authenticating user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
