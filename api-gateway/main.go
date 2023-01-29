package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iamvasanth07/showcase/api-gateway/config"
	"github.com/iamvasanth07/showcase/api-gateway/routes"
	userConfig "github.com/iamvasanth07/showcase/user/config"
	videoConfig "github.com/iamvasanth07/showcase/video/config"
)

// serve the user routes using gin

func main() {
	settings := config.GetSettings()
	userRoutes := routes.NewUserRoutes(userConfig.GetSettings())
	videoRoutes := routes.NewVideoRoutes(videoConfig.GetSettings())
	r := gin.Default()
	userRoutes.RegisterUserSvcRoutes(r)
	videoRoutes.RegisterVideoSvcRoutes(r)
	r.Run(fmt.Sprintf("%s:%s", settings.Server.HTTPHost, settings.Server.HTTPPort))
}
