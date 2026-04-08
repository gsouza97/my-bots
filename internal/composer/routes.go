package composer

import (
	"github.com/gin-gonic/gin"
	"github.com/gsouza97/my-bots/internal/interfaces/http/handlers"
	"github.com/gsouza97/my-bots/internal/interfaces/http/middleware"
	"github.com/gsouza97/my-bots/internal/interfaces/http/routes"
)

type RoutesComposer struct {
	AlertsRoutes *routes.AlertsRoutes
	LoginRoutes  *routes.LoginRoutes
	LoansRoutes  *routes.LoansRoutes
	PoolsRoutes  *routes.PoolsRoutes
	BillsRoutes  *routes.BillsRoutes
}

type HandlerComposer struct {
	AlertsHandler *handlers.AlertsHandler
	LoginHandler  *handlers.LoginHandler
	LoansHandler  *handlers.LoansHandler
	PoolsHandler  *handlers.PoolsHandler
	BillsHandler  *handlers.BillsHandler
}

func NewHandlerComposer(services *ServicesComposer) *HandlerComposer {
	return &HandlerComposer{
		AlertsHandler: handlers.NewAlertsHandler(services.AlertsService),
		LoginHandler:  handlers.NewLoginHandler(services.AuthService),
		LoansHandler:  handlers.NewLoansHandler(services.LoansService),
		PoolsHandler:  handlers.NewPoolsHandler(services.PoolsService),
		BillsHandler:  handlers.NewBillsHandler(services.BillsService),
	}
}

func NewRoutesComposer(handlers *HandlerComposer) *RoutesComposer {
	return &RoutesComposer{
		AlertsRoutes: routes.NewAlertsRoutes(handlers.AlertsHandler),
		LoginRoutes:  routes.NewLoginRoutes(handlers.LoginHandler),
		LoansRoutes:  routes.NewLoansRoutes(handlers.LoansHandler),
		PoolsRoutes:  routes.NewPoolsRoutes(handlers.PoolsHandler),
		BillsRoutes:  routes.NewBillsRoutes(handlers.BillsHandler),
	}
}

func (rc *RoutesComposer) RegisterRoutes(engine *gin.Engine, userToken string) {
	authMiddleware := middleware.AuthMiddleware(userToken)

	// Rotas públicas
	routes.StartHealthRoutes(engine)
	rc.LoginRoutes.StartLoginRoutes(engine)

	// Rotas protegidas
	rc.AlertsRoutes.StartAlertsRoutes(engine, authMiddleware)
	rc.LoansRoutes.StartLoansRoutes(engine, authMiddleware)
	rc.PoolsRoutes.StartPoolsRoutes(engine, authMiddleware)
	rc.BillsRoutes.StartBillsRoutes(engine, authMiddleware)
}
