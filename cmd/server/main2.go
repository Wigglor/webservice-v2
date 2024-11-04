package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/Wigglor/webservice-v2/handlers"
// 	"github.com/Wigglor/webservice-v2/repository"
// 	"github.com/Wigglor/webservice-v2/repository/database"
// 	"github.com/Wigglor/webservice-v2/router"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/joho/godotenv"
// )

// type Config struct {
// 	DB     database.Config
// 	Server ServerConfig
// }

// type ServerConfig struct {
// 	Addr string
// }

// // Main function remains in the main file
// func main2() {
// 	cfg, err := loadConfig()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	app, err := initializeApp(cfg)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize application: %v", err)
// 	}
// 	defer app.Close()

// 	if err := app.Run(); err != nil {
// 		log.Fatalf("Application error: %v", err)
// 	}
// }

// // Application struct and methods remain in main file
// type App struct {
// 	Config Config
// 	DB     *pgxpool.Pool
// 	Router http.Handler
// 	Server *http.Server
// }

// func initializeApp(cfg Config) (*App, error) {
// 	// Initialize database connection
// 	pool, err := database.ConnectDB(cfg.DB)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Set up repositories, handlers, and router
// 	userRepo := repository.NewUserRepository(pool)
// 	userHandler := handlers.NewUserHandler(userRepo)
// 	r := router.Routes(userHandler)

// 	// Create the server
// 	srv := &http.Server{
// 		Addr:    cfg.Server.Addr,
// 		Handler: r,
// 	}

// 	return &App{
// 		Config: cfg,
// 		DB:     pool,
// 		Router: r,
// 		Server: srv,
// 	}, nil
// }

// func (app *App) Run() error {
// 	// Handle graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// 	// Start the server in a goroutine
// 	go func() {
// 		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("HTTP server error: %v", err)
// 		}
// 	}()
// 	log.Printf("Server Started on %s", app.Server.Addr)

// 	<-quit
// 	log.Print("Server Stopped")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := app.Server.Shutdown(ctx); err != nil {
// 		return err
// 	}

// 	log.Print("Server Exited Properly")
// 	return nil
// }

// func (app *App) Close() {
// 	if app.DB != nil {
// 		app.DB.Close()
// 	}
// }

// // Configuration loading function remains in main file
// func loadConfig2() (Config, error) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Println("failed godotenv.Load")
// 		return Config{}, fmt.Errorf("error loading .env file: %w", err)
// 	}

// 	dbURL := ConcatDSN()
// 	fmt.Println("dbURL: ", dbURL)

// 	dbConfig := database.Config{
// 		DSN:             dbURL,
// 		MaxConns:        10,
// 		MinConns:        5,
// 		MaxConnLifetime: time.Hour,
// 		MaxConnIdleTime: 30 * time.Minute,
// 	}

// 	serverConfig := ServerConfig{
// 		Addr: ":8080",
// 	}

// 	return Config{
// 		DB:     dbConfig,
// 		Server: serverConfig,
// 	}, nil
// }

// func ConcatDSN2() string {
// 	host := os.Getenv("DB_HOST")
// 	username := os.Getenv("DB_USER")
// 	password := os.Getenv("DB_PASSWORD")
// 	databaseName := os.Getenv("DB_NAME")
// 	port := os.Getenv("DB_PORT")

// 	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, databaseName)
// }
