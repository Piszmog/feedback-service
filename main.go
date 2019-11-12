package main

import (
	"github.com/Piszmog/feedback-service/db"
	"github.com/Piszmog/feedback-service/transport"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	defaultPort           = "8080"
	environmentPort       = "PORT"
	environmentDBDatabase = "DB_DATABASE"
	environmentDBHost     = "DB_HOST"
	environmentDBPassword = "DB_PASSWORD"
	environmentDBPort     = "DB_PORT"
	environmentDBUsername = "DB_USERNAME"
)

func main() {
	//
	// Connect to the DB
	//
	mysql := connectToDB()
	defer mysql.Close()
	//
	// Create the table if does not exist
	//
	if err := mysql.CreateFeedbackTableIfNotExists(); err != nil {
		log.Println(err)
		return
	}
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
		DB:           mysql,
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

func connectToDB() *db.MySQL {
	//
	// Get env variable for the DB
	//
	username := os.Getenv(environmentDBUsername)
	password := os.Getenv(environmentDBPassword)
	host := os.Getenv(environmentDBHost)
	dbPort := os.Getenv(environmentDBPort)
	database := os.Getenv(environmentDBDatabase)
	options := db.Options{
		Username:     username,
		Password:     password,
		Host:         host,
		Port:         dbPort,
		DatabaseName: database,
		ParseTime:    true,
	}
	mysql, err := db.Open(options)
	if err != nil {
		log.Fatalln(err)
	}
	return mysql
}

func gracefulShutdown(srv *transport.HTTPServer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	srv.Shutdown(5 * time.Second)
	log.Println("shutting down...")
	os.Exit(0)
}
