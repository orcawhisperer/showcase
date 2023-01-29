package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iamvasanth07/showcase/api-gateway/config"
	"github.com/iamvasanth07/showcase/api-gateway/routes"
	userConfig "github.com/iamvasanth07/showcase/user/config"
)

// serve the user routes using gin

func main() {
	settings := config.GetSettings()
	userRoutes := routes.NewUserRoutes(userConfig.GetSettings())

	r := gin.Default()
	r = r.Group("/api/v1")
	userRoutes.RegisterUserSvcRoutes(r)
	r.Run(settings.Server.HTTPPort)
}
