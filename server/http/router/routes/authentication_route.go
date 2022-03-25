package routes

import (
	"svc-whatsapp/server/http/handlers"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/server/http/middlewares"

	"github.com/gin-gonic/gin"
)

type AuthenticationRoute struct {
	RouteGroup *gin.RouterGroup
	Handler    handler.Handler
}

func NewAuthenticationRoute(routeGroup *gin.RouterGroup, handler handler.Handler) AuthenticationRoute {
	return AuthenticationRoute{RouteGroup: routeGroup, Handler: handler}
}

func (route AuthenticationRoute) RegisterRoute() {

	//middleware
	middlewareJwt := middlewares.NewJWTMiddleware(route.Handler.Contract)

	// handlers
	authenticationHandler := handlers.NewAuthenticationHandler(route.Handler)

	// Users Route
	route.RouteGroup.POST("/login", authenticationHandler.Login)
	route.RouteGroup.POST("/register", authenticationHandler.Register)

	// using jwt
	loginOnly := route.RouteGroup.Use(middlewareJwt.LoginOnly)
	loginOnly.GET("/jwt", authenticationHandler.GetJwt)

}
