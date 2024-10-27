package main

import (
	"context"
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
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	/*dbConfig, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	config, err := pgxpool.ParseConfig(dbConfig.DSN)
	if err != nil {
		log.Fatalf("Failed to parse database configuration: %v", err)
		return
	}

	config.MaxConns = dbConfig.MaxConns
	config.MinConns = dbConfig.MinConns
	config.MaxConnLifetime = dbConfig.MaxConnLifetime
	config.MaxConnIdleTime = dbConfig.MaxConnIdleTime*/

	/*// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	defer func() {
		// extra handling here
		cancel()
	}()*/

	/*pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}*/
	pool, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer pool.Close()

	/*err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to Ping...: %v", err)
		pool.Close() // should i have this here???
		return
	}*/

	var wg sync.WaitGroup

	userRepo := repository.NewUserRepository(pool)
	userHandler := handlers.NewUserHandler(userRepo) // changfrom router to controller/handler folder
	// userHandler := router.NewUserHandler(userRepo) // changfrom router to controller/handler folder
	router := router.Routes(userHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	wg.Add(1) // Does nothing atm. implement in handlers if neccessary.
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// if err := srv.ListenAndServe();  err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
			// log.Printf("HTTP server error: %v", err)
			// quit <- os.Interrupt
		}
		// // if err := http.ListenAndServe(":8080", router); err != nil {
		// if err := srv.ListenAndServe(); err != nil {
		// 	log.Fatalf("HTTP server error: %v", err)
		// }
	}()
	log.Print("Server Started")

	<-quit
	log.Print("Server Stopped")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	//why shutdown context???

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	wg.Wait()
	log.Print("All goroutines have finished")
	log.Print("Server Exited Properly")
}

/*func ConcatDSN() string {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("%s://%s:%s@db:%s/%s", host, username, password, port, databaseName)
	// return fmt.Sprintf("%s://%s:%s@localhost:%s/%s", host, username, password, port, databaseName)
}

type Config struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func loadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed godotenv.Load")
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	dbURL := ConcatDSN()
	fmt.Println("dbURL: ", dbURL)

	return Config{
		DSN:             dbURL,
		MaxConns:        10,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}, nil
}*/
