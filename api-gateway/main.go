package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/iamvasanth07/showcase/api-gateway/config"
	usrConfig "github.com/iamvasanth07/showcase/user/config"
	userSvc "github.com/iamvasanth07/showcase/user/service"
	videoSvc "github.com/iamvasanth07/showcase/video/service"
)

func main() {

	settings := config.GetSettings()

	logger := log.New(os.Stdout, "api-gateway", log.LstdFlags)

	usrSettings := usrConfig.GetSettings()

	logger.Println("user settings loaded: ", usrSettings.Database.Host, usrSettings.Database.Port, usrSettings.Database.User, usrSettings.Database.Password, usrSettings.Database.Name)

	// get user service http handler
	userHandler := userSvc.GetHTTPHandler()
	// get video service http handler
	videoHandler := videoSvc.GetHTTPHandler()

	// create a new serve mux and register handlers
	mux := http.NewServeMux()
	mux.Handle("/v1/user", userHandler)
	mux.Handle("/api/v1/videos", videoHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", settings.Server.HTTPHost, settings.Server.HTTPPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logger.Println("HTTP Server started on port: " + settings.Server.HTTPPort)
	if err := http.Serve(lis, mux); err != nil {
		log.Fatalf("failed to serv http server: %v", err)
	}

}
