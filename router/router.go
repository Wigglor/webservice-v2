package router

import (
	"fmt"
	"net/http"

	"github.com/Wigglor/webservice-v2/handlers"
	"github.com/Wigglor/webservice-v2/middlewares"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// func Routes(handler *UserHandler) http.Handler {
func Routes(handler *handlers.UserHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(chi_middleware.Recoverer)
	router.Use(chi_middleware.Logger)
	router.Use(middlewares.LoggerMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api", func(r chi.Router) {
		r.Get("/users", handler.GetUsers)
		r.Get("/user/{id}", handler.GetUserById)

		r.Route("/auth", func(r chi.Router) {
			r.With(middlewares.EnsureValidToken()).Get("/protected", helloAuth)
			// r.Use(middlewares.EnsureValidToken())
			r.Get("/", helloWorld) // DELETE /articles/123
		})

	})

	return router
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, World!")
}

func helloAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, World!")
}

/*type UserHandler struct {
	Repo repository.UserRepository
	// wg   *sync.WaitGroup
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, World!")
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.QueryAllUsers(r.Context())
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
	user, err := h.Repo.GetUserByID(r.Context(), int32(userId))
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}*/
