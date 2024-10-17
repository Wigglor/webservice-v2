package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/Wigglor/webservice-v2/repository"
	"github.com/Wigglor/webservice-v2/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	dbConfig, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("after loadConfig")
	//dbpool, err := pgxpool.New(context.Background(), dbConfig.DSN)

	config, err := pgxpool.ParseConfig(dbConfig.DSN)
	if err != nil {
		log.Fatalf("Failed to parse database configuration: %v", err)
		return
		// log.Fatalf("Failed to parse database configuration: %v", err)
	}
	fmt.Println("after ParseConfig")

	// Configure the pool settings
	config.MaxConns = dbConfig.MaxConns
	config.MinConns = dbConfig.MinConns
	config.MaxConnLifetime = dbConfig.MaxConnLifetime
	config.MaxConnIdleTime = dbConfig.MaxConnIdleTime

	// Create the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	fmt.Println("after NewWithConfig")

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to Ping...: %v", err)
		pool.Close()
		return
		// log.Fatalf("Failed to test the connection: %v", err)
	}
	fmt.Println("after Ping")

	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	fmt.Println("after NewUserRepository", userRepo)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	userHandler := router.NewUserHandler(userRepo)
	fmt.Println("after NewUserHandler")
	router := router.Routes(userHandler)
	/*router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api", func(r chi.Router) {
		r.Get("/users", userHandler.GetUsers)        // attaching a method to the App struct -> pointer receiver
		r.Get("/user/{id}", userHandler.GetUserById) // attaching a method to the App struct -> pointer receiver
		r.Get("/", helloWorld)
		// r.Get("/user/{id}", controller.GetUserById) // attaching a method to the App struct -> pointer receiver
		// r.Get("/users", controller.GetUser)         // attaching a method to the App struct -> pointer receiver
		// r.Get("/user/{id}", controller.GetUserById) // attaching a method to the App struct -> pointer receiver

	})*/

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

// type UserHandler struct {
// 	Repo repository.UserRepository
// 	// Context context.Context
// }

// func NewUserHandler(repo repository.UserRepository) *UserHandler {
// 	return &UserHandler{Repo: repo}
// }

/*func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.QueryAllUsers()
	if err != nil {
		log.Fatalf("QueryAllUsers error: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "id")
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// user, err := h.Repo.GetUserByID(r.Context(), int32(userId))
	user, err := h.Repo.GetUserByID(int32(userId))
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}*/

// func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	users, err := h.Repo.QueryCreateUser()
// 	if err != nil {
// 		log.Fatalf("QueryAllUsers error: %v", err)
// 		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(users); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 	}

// }

/*func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, World!")
}*/

func ConcatDSN() string {
	host := "postgresql"
	username := "webservice_dev_user"
	password := "yourpassword"
	databaseName := "webservice_dev"
	port := "5432"
	// host := os.Getenv("DB_HOST")
	// username := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// databaseName := os.Getenv("DB_NAME")
	// port := os.Getenv("DB_PORT")

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
}
