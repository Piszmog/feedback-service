package main

import (
	"github.com/Piszmog/feedback-service/transport"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	defaultPort     = "8080"
	environmentPort = "PORT"
)

func main() {
	//
	// Connect to the DB
	//
	// TODO
	//
	// Get the port
	//
	port := os.Getenv(environmentPort)
	if len(port) == 0 {
		port = defaultPort
	}
	//
	// Create the HTTP server and run it
	//
	srv := &transport.HTTPServer{
		Port:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		//DB: // TODO
	}
	go func() {
		if err := srv.Start(); err != nil {
			log.Println(err)
		}
	}()
	//
	// If any shutdown signals come, then try to gracefully shut the server down
	//
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *transport.HTTPServer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	srv.Shutdown(5 * time.Second)
	log.Println("shutting down...")
	os.Exit(0)
}
