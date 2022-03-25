package routes

import (
	"svc-whatsapp/server/http/handlers"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/server/http/middlewares"

	"github.com/gin-gonic/gin"
)

type DevicesRoute struct {
	RouteGroup *gin.RouterGroup
	Handler    handler.Handler
}

func NewDevicesRoute(routeGroup *gin.RouterGroup, handler handler.Handler) DevicesRoute {
	return DevicesRoute{RouteGroup: routeGroup, Handler: handler}
}

func (route DevicesRoute) RegisterRoute() {

	//middleware
	middlewareJwt := middlewares.NewJWTMiddleware(route.Handler.Contract)

	// handlers
	DevicesHandler := handlers.NewDeviceHandler(route.Handler)

	// login only
	loginOnly := route.RouteGroup.Use(middlewareJwt.LoginOnly)
	loginOnly.GET("/devices/filter", DevicesHandler.Filter)

}
