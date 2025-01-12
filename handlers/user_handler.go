package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"strings"

	"github.com/Wigglor/webservice-v2/repository"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	Repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
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
	userIdStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	println(userIdStr)
	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	println(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return

	}
	user, err := h.Repo.GetUserByID(r.Context(), int32(userId))
	if err != nil {
		log.Printf("Failed to fetch user with ID %d: %v", userId, err)
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetOrCreateUserBySubId(w http.ResponseWriter, r *http.Request) {
	// type RequestBody struct {
	// 	Name               string `json:"name"`
	// 	Email              string `json:"email"`
	// 	VerificationStatus bool   `json:"verificationStatus"`
	// 	SubId              string `json:"subId"`
	// }

	var reqBody repository.CreateUserParams

	// 2. Decode the JSON from the request body into that struct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		// Handle error (e.g., malformed JSON)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// userSubIdStr := strings.TrimPrefix(r.URL.Path, "/api/check-user/")
	// println(userSubIdStr)
	// userId, err := strconv.ParseInt(userSubIdStr, 10, 32)
	// println(userId)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return

	// }
	user, err := h.Repo.CheckUserBySubId(r.Context(), reqBody.SubId)
	if err != nil {
		log.Printf("Failed to fetch user with Sub ID %s: %v", reqBody.SubId, err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No rows error")
			log.Println(pgx.ErrNoRows)

			createdUser, err := h.Repo.QueryCreateUser(r.Context(), reqBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println(createdUser)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(createdUser); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}
		// http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		// return
		return
	}
	// log.Println(user)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		// if err := json.NewEncoder(w).Encode(reqBody); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var reqBody repository.CreateOrganizationParams

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		// Handle error (e.g., malformed JSON)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrgUser, err := h.Repo.QueryCreateOrganization(r.Context(), reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdOrgUser); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
