package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Wigglor/webservice-v2/handlers"
	"github.com/Wigglor/webservice-v2/repository"
	"github.com/Wigglor/webservice-v2/repository/database"
	"github.com/Wigglor/webservice-v2/router"
	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main2() {

	dbConfig, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := database.ConnectDB(dbConfig) // test this with testcontainers
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer pool.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	/*userRepo := repository.NewUserRepository(pool)
	userHandler := handlers.NewUserHandler(userRepo)
	router := router.Routes(userHandler)*/

	router := setupRouter2(pool)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	log.Print("Server Started")

	<-quit
	log.Print("Server Stopped")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	wg.Wait()
	log.Print("All goroutines have finished")
	log.Print("Server Exited Properly")
}

func setupRouter2(pool *pgxpool.Pool) http.Handler {
	// Initialize the repository with the database connection pool
	userRepo := repository.NewUserRepository(pool)

	// Create the user handler with the repository
	userHandler := handlers.NewUserHandler(userRepo)

	// Set up the routes and return the router
	return router.Routes(userHandler)
}

func ConcatDSN2() string {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, databaseName)
}

func loadConfig2() (database.Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed godotenv.Load")
		return database.Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	dbURL := ConcatDSN()
	fmt.Println("dbURL: ", dbURL)

	return database.Config{
		DSN:             dbURL,
		MaxConns:        10,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}, nil
}
