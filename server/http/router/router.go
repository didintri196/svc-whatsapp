package router

import (
	"net/http"
	handler "svc-whatsapp/server/http/handlers/helper"
	"svc-whatsapp/server/http/middlewares"
	"svc-whatsapp/server/http/router/routes"
	"time"

	"github.com/gin-contrib/cors"

	statsApi "github.com/appleboy/gin-status-api"
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler.Handler
}

func NewRouter(handler handler.Handler) Router {
	return Router{handler}
}

func (router Router) RegisterRoutes() {

	app := router.Contract.App
	socket := router.Contract.Socket
	queryMiddleware := middlewares.NewQueryiddleware(router.Contract)

	// middleware
	app.Use(gin.Recovery())

	//CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PUT", "DELETE", "PATCH", "WSS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Accept-Language"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Root Group
	rootGroup := app.Group("/svc-whatsapp")

	// Route for check health
	rootGroup.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "svc-whatsapp ok!")
	})

	// Route for check stats
	rootGroup.GET("/stats", statsApi.GinHandler)

	// Socket io init
	app.GET("/socket.io/*any", queryMiddleware.TokenisFilled, gin.WrapH(socket))
	app.POST("/socket.io/*any", queryMiddleware.TokenisFilled, gin.WrapH(socket))
	app.StaticFS("/public", http.Dir("../../domain/assets"))

	// socket middleware

	// Qrcode route
	whatsappSocketRout := routes.NewWhatsappSocketioRoute(socket, router.Handler)
	whatsappSocketRout.RegisterRoute()

	// Whatsapp route
	whatsappRoute := routes.NewWhatsappRoute(rootGroup, router.Handler)
	whatsappRoute.RegisterRoute()

	// Authentication route
	authenticationRoute := routes.NewAuthenticationRoute(rootGroup, router.Handler)
	authenticationRoute.RegisterRoute()

	// Devices route
	DevicecsRoute := routes.NewDevicesRoute(rootGroup, router.Handler)
	DevicecsRoute.RegisterRoute()

}
