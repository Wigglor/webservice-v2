package router

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Wigglor/webservice-v2/handlers"
	"github.com/Wigglor/webservice-v2/middlewares"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	// "github.com/go-chi/chi/v5"
	// chi_middleware "github.com/go-chi/chi/v5/middleware"
	// "github.com/go-chi/cors"
)

// func Routes(handler *UserHandler) http.Handler {
func Routes(handler *handlers.UserHandler) http.Handler {

	// finalHandler := http.HandlerFunc(helloAuth)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/users", handler.GetUsers)
	mux.HandleFunc("GET /api/user/{id}", handler.GetUserById)
	// mux.HandleFunc("GET /api/auth/protected", middlewares.EnsureValidToken(finalHandler))
	// This route is only accessible if the user has a valid access_token.
	mux.Handle("/api/private", middlewares.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS Headers.
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello from a private endpoint! You need to be authenticated to see this message."}`))
		}),
	))
	mux.Handle("/api/private2", ValidateJWT(http.HandlerFunc(helloAuth)))
	mux.Handle("/api/private3", middlewares.EnsureValidToken()(http.HandlerFunc(helloAuth)))
	return mux
	/*router := chi.NewRouter()
	router.Use(chi_middleware.Recoverer)
	router.Use(chi_middleware.Logger)
	router.Use(ValidateJWT)
	router.Use(middlewares.LoggerMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access granted to protected route!"))
	})

	router.Route("/api", func(r chi.Router) {
		r.Get("/users", handler.GetUsers)
		r.Get("/user/{id}", handler.GetUserById)
		r.With().Get("/protected", helloAuth)

		r.Route("/auth", func(r chi.Router) {
			// r.Use(middlewares.EnsureValidToken)
			// r.Use(middlewares.ValidateJWT)
			//r.With(middlewares.EnsureValidToken).Get("/protected", helloAuth)
			// r.Use(middlewares.EnsureValidToken())
			r.Get("/", helloWorld) // DELETE /articles/123
		})

	})

	return router*/
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there, World!")
}

func helloAuth(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello there, World!")
	w.Write([]byte("OK"))
}

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
		if err != nil {
			log.Fatalf("Failed to parse the issuer url: %v", err)
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{os.Getenv("AUTH0_AUDIENCE")},
		)
		if err != nil {
			log.Fatalf("Failed to set up the jwt validator")
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Encountered error while validating JWT: %v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"Failed to validate JWT."}`))
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)

		// Log authorization header for debugging
		authHeader := r.Header.Get("Authorization")
		log.Printf("Authorization Header: %s", authHeader)

		// Pass through the middleware, which may modify the request context
		middleware.CheckJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Token validated successfully.")
			next.ServeHTTP(w, r)
		})).ServeHTTP(w, r)

		middleware.CheckJWT(next).ServeHTTP(w, r)
	})
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
