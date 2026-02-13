package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/handlers"
)

type PoolsRoutes struct {
	poolsHandler *handlers.PoolsHandler
}

func NewPoolsRoutes(poolsHandler *handlers.PoolsHandler) *PoolsRoutes {
	return &PoolsRoutes{
		poolsHandler: poolsHandler,
	}
}

func (rt *PoolsRoutes) StartPoolsRoutes(r *gin.Engine, middleware ...gin.HandlerFunc) {
	r.GET("/pools", append(middleware, rt.poolsHandler.GetAllPools)...)
}
