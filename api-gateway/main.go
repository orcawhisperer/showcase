// api gateway for the grpc microservices user and video with all the CRUD operations

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// create a new server mux and register the handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler)
	mux.HandleFunc("/videos", videoHandler)

	// create a new server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start the server
	go func() {
		log.Println("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting server: %v ", err)
		}
	}()

	// create a channel to listen for an interrupt or terminate signal from the OS
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// block until a signal is received
	<-c

	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	s.Shutdown(ctx)

	// Optionally, you could run s.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.

	log.Println("shutting down")
	os.Exit(0)
}

//
