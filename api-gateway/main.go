package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/iamvasanth07/showcase/api-gateway/config"
	userSvc "github.com/iamvasanth07/showcase/user/service"
)

func main() {

	settings := config.GetSettings()

	logger := log.New(os.Stdout, "api-gateway", log.LstdFlags)

	// get user service http handler
	userHandler := userSvc.GetHTTPHandler()

	// create a new serve mux and register handlers
	mux := http.NewServeMux()
	mux.Handle("/user/", userHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", settings.Server.HTTPHost, settings.Server.HTTPPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logger.Println("HTTP Server started on port: " + settings.Server.HTTPPort)
	if err := http.Serve(lis, mux); err != nil {
		log.Fatalf("failed to serv http server: %v", err)
	}

}
