package routes

import (
	"svc-whatsapp/server/http/handlers"
	handler "svc-whatsapp/server/http/handlers/helper"

	"github.com/gin-gonic/gin"
)

type WhatsappRoute struct {
	RouteGroup *gin.RouterGroup
	Handler    handler.Handler
}

func NewWhatsappRoute(routeGroup *gin.RouterGroup, handler handler.Handler) WhatsappRoute {
	return WhatsappRoute{RouteGroup: routeGroup, Handler: handler}
}

func (route WhatsappRoute) RegisterRoute() {

	// handlers
	whatsappHandler := handlers.NewWhatsappApiHandler(route.Handler)

	// Test Produce
	route.RouteGroup.POST("/produce", whatsappHandler.TestProduce)
}
