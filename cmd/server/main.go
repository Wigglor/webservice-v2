package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// "time"

	"github.com/Wigglor/webservice-v2/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := ConcatDSN()
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbpool.Close()

	userRepo := repository.NewUserRepository(dbpool)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	userHandler := NewUserHandler(userRepo)

	router := chi.NewRouter()
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
		r.Get("/usersss", userHandler.GetUsers) // attaching a method to the App struct -> pointer receiver
		r.Get("/", helloWorld)
		// r.Get("/user/{id}", controller.GetUserById) // attaching a method to the App struct -> pointer receiver
		// r.Get("/users", controller.GetUser)         // attaching a method to the App struct -> pointer receiver
		// r.Get("/user/{id}", controller.GetUserById) // attaching a method to the App struct -> pointer receiver

	})
}

type UserHandler struct {
	Repo    repository.UserRepository
	Context context.Context
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.QueryAllUsers()

	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, World!")
}

func ConcatDSN() string {
	// host := "postgresql"
	// username := "webservice_dev_user"
	// password := "yourpassword"
	// databaseName := "webservice_dev"
	// port := "5432"
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("%s://%s:%s@db:%s/%s", host, username, password, port, databaseName)
}
