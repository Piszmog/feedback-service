package main

import (
	"database/sql"
	"fmt"
	"github.com/Piszmog/feedback-service/db"
	"github.com/Piszmog/feedback-service/transport"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	defaultDatabase       = "ubisoft"
	defaultDBHost         = "localhost"
	defaultDBPort         = "3306"
	defaultHost           = "localhost"
	defaultPort           = "8080"
	environmentHost       = "HOST"
	environmentPort       = "PORT"
	environmentDBDatabase = "DB_DATABASE"
	environmentDBHost     = "DB_HOST"
	environmentDBPassword = "DB_PASSWORD"
	environmentDBPort     = "DB_PORT"
	environmentDBUsername = "DB_USERNAME"
)

func main() {
	start := time.Now()
	log.Println("Starting application...")
	//
	// Connect to the DB
	//
	mysql, err := createMySQLDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer mysql.Close()
	//
	// Create the table if does not exist
	//
	if err := mysql.CreateFeedbackTableIfNotExists(); err != nil {
		log.Println(err)
		return
	}
	//
	// Get the host and port
	//
	host := os.Getenv(environmentHost)
	if len(host) == 0 {
		host = defaultHost
	}
	port := os.Getenv(environmentPort)
	if len(port) == 0 {
		port = defaultPort
	}
	//
	// Create the HTTP server and run it
	//
	srv := &transport.HTTPServer{
		Host:         host,
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
	log.Printf("Application started in %f seconds\n", time.Since(start).Seconds())
	//
	// If any shutdown signals come, then try to gracefully shut the server down
	//
	gracefulShutdown(srv)
}

func createMySQLDB() (*db.MySQL, error) {
	//
	// Get env variable for the DB
	//
	username := os.Getenv(environmentDBUsername)
	password := os.Getenv(environmentDBPassword)
	host := os.Getenv(environmentDBHost)
	if len(host) == 0 {
		log.Println("Defaulting to default DB host name 'localhost'")
		host = defaultDBHost
	}
	dbPort := os.Getenv(environmentDBPort)
	if len(dbPort) == 0 {
		log.Println("Defaulting to default DB port '3306'")
		dbPort = defaultDBPort
	}
	database := os.Getenv(environmentDBDatabase)
	if len(database) == 0 {
		log.Println("Defaulting to default database name 'ubisoft'")
		database = defaultDatabase
	}
	options := db.Options{
		Username:     username,
		Password:     password,
		Host:         host,
		Port:         dbPort,
		DatabaseName: database,
		ParseTime:    true,
	}
	//
	// Validate the provided options
	//
	if err := options.Validate(); err != nil {
		return nil, err
	}
	//
	// Connect to the DB
	//
	dbConnection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		options.Username, options.Password, options.Host, options.Port, options.DatabaseName, options.ParseTime))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the MySQL DB: %w", err)
	}
	//
	// Ensure we can talk to the DB
	//
	if err := dbConnection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}
	log.Printf("Successfully connected to %s database\n", database)
	return &db.MySQL{DB: dbConnection}, nil
}

func gracefulShutdown(srv *transport.HTTPServer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	srv.Shutdown(5 * time.Second)
	log.Println("shutting down...")
	os.Exit(0)
}
