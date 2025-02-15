package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"strings"

	"github.com/Wigglor/webservice-v2/repository"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
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

	claims, ok := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		http.Error(w, "Failed to get validated token from context", http.StatusUnauthorized)
		return
	}

	var reqBody repository.CreateUserParams

	// 2. Decode the JSON from the request body into that struct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		// Handle error (e.g., malformed JSON)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Repo.CheckUserBySubId(r.Context(), claims.RegisteredClaims.Subject)
	subId := claims.RegisteredClaims.Subject
	if err != nil {
		log.Printf("Failed to fetch user with Sub ID %s: %v", claims.RegisteredClaims.Subject, err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println(pgx.ErrNoRows)

			log.Println("----------------------------------------------------")
			log.Println("----------------------------------------------------")
			log.Println(reqBody)
			log.Println(&reqBody)
			log.Println("----------------------------------------------------")
			log.Println("----------------------------------------------------")
			createdUser, err := h.Repo.QueryCreateUser(r.Context(), reqBody, subId)
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
		log.Printf("Failed to fetch user with Sub ID %s: %v", claims.RegisteredClaims.Subject, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// In case the user signs in before verifying the email
	println("checking VerificationStatus...")
	if !user.VerificationStatus && reqBody.VerificationStatus {
		println("setting VerificationStatus to true...")
		updatedVerificationStatus, err := h.Repo.UpdateVerificationStatus(r.Context(), claims.RegisteredClaims.Subject)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.VerificationStatus = updatedVerificationStatus
	}

	if user.SetupStatus == "completed" {
		println("status is completed")
		userOrg, err := h.Repo.QueryOrganization(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(userOrg)
		type CombinedData struct {
			User    repository.User           `json:"user"`
			UserOrg []repository.Organization `json:"userOrg"`
		}
		combinedData := CombinedData{
			User:    user,
			UserOrg: userOrg,
		}
		fmt.Println(userOrg)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(combinedData); err != nil {
			// if err := json.NewEncoder(w).Encode(reqBody); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	userPayload := map[string]interface{}{
		"user": user,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userPayload); err != nil {
		// if err := json.NewEncoder(w).Encode(reqBody); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CreateUserForOrg(w http.ResponseWriter, r *http.Request) {
	/*
		1. get user that makes the call based on auth0 id
		2. check if user is admin
			-	if yes, next step
			-	if no, return error
		3. check if user is part of the org in question (get org id from payload)
		4. check how many users are allowed for that org and how many there currently are
		5. create user in in db, auth0 and link it to the organization
		6. send invite to user email via SES or similar


	*/
	// var reqBody repository.CreateOrganizationParams

	// err := json.NewDecoder(r.Body).Decode(&reqBody)
	// if err != nil {
	// 	// Handle error (e.g., malformed JSON)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// createdOrgUser, err := h.Repo.QueryCreateOrganization(r.Context(), reqBody)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(createdOrgUser); err != nil {
	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// }
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
