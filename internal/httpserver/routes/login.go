package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/handlers"
)

type LoginRoutes struct {
	loginHandler *handlers.LoginHandler
}

func NewLoginRoutes(loginHandler *handlers.LoginHandler) *LoginRoutes {
	return &LoginRoutes{
		loginHandler: loginHandler,
	}
}

func (rt *LoginRoutes) StartLoginRoutes(r *gin.Engine) {
	r.POST("/login", rt.loginHandler.Login)
}
