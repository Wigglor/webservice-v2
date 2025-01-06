package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Wigglor/webservice-v2/handlers"
	"github.com/Wigglor/webservice-v2/repository"

	// "github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

type mockUserModel struct{}

func (m *mockUserModel) QueryAllUsers(ctx context.Context) ([]repository.User, error) {
	var users []repository.User

	users = append(users, repository.User{ID: 1, Name: "Joe Doe test", Email: "johndoe@email.com", SubID: "subid_123abc", VerificationStatus: true, SetupStatus: "pending"})
	users = append(users, repository.User{ID: 2, Name: "Jane Doe test", Email: "janedoe@email.com", SubID: "subid_456def", VerificationStatus: true, SetupStatus: "pending"})

	return users, nil
}

func (m *mockUserModel) GetUserByID(ctx context.Context, id int32) (repository.User, error) {
	user := repository.User{ID: 1, Name: "Joe Doe", Email: "johndoe@email.com", SubID: "subid_123abc", VerificationStatus: true, SetupStatus: "pending"}
	return user, nil
}

func (m *mockUserModel) CheckUserBySubId(ctx context.Context, subId string) (repository.User, error) {
	user := repository.User{ID: 1, Name: "Joe Doe", Email: "johndoe@email.com", SubID: "subid_123abc", VerificationStatus: true, SetupStatus: "pending"}
	return user, nil
}

func (m *mockUserModel) QueryCreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error) {
	user := repository.User{ID: 1, Name: "Joe Doe", Email: "johndoe@email.com", SubID: "subid_123abc", VerificationStatus: true, SetupStatus: "pending"}
	return user, nil
}
func TestGetUsers(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	/*handler := router.UserHandler{Repo: &mockUserModel{}}*/
	handler := handlers.UserHandler{Repo: &mockUserModel{}}
	http.HandlerFunc(handler.GetUsers).ServeHTTP(rec, req)
	expected := []repository.User{
		{
			ID:                 1,
			Name:               "Joe Doe test",
			Email:              "johndoe@email.com",
			SubID:              "subid_123abc",
			VerificationStatus: true,
			SetupStatus:        "pending",
			CreatedAt:          pgtype.Timestamptz{Valid: false},
			UpdatedAt:          pgtype.Timestamptz{Valid: false},
		},
		{
			ID:                 2,
			Name:               "Jane Doe test",
			Email:              "janedoe@email.com",
			SubID:              "subid_456def",
			VerificationStatus: true,
			SetupStatus:        "pending",
			CreatedAt:          pgtype.Timestamptz{Valid: false},
			UpdatedAt:          pgtype.Timestamptz{Valid: false},
		},
	}
	var obtained []repository.User
	err := json.Unmarshal(rec.Body.Bytes(), &obtained)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Compare the expected and obtained data
	if !reflect.DeepEqual(expected, obtained) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}

	// Alternatively, using testify's assert
	assert.Equal(t, expected, obtained, "The expected and obtained users should be equal")
	// if expected != rec.Body.String() {
	// 		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	// 	}
}

func TestGetUserById(t *testing.T) {
	// Create the mock repository and handler
	handler := handlers.UserHandler{Repo: &mockUserModel{}}

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/api/user/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rec := httptest.NewRecorder()

	// Call the handler directly
	http.HandlerFunc(handler.GetUserById).ServeHTTP(rec, req)

	// Define the expected user data
	expected := repository.User{
		ID:                 1,
		Name:               "Joe Doe",
		Email:              "johndoe@email.com",
		SubID:              "subid_123abc",
		VerificationStatus: true,
		SetupStatus:        "pending",
		CreatedAt:          pgtype.Timestamptz{Valid: false},
		UpdatedAt:          pgtype.Timestamptz{Valid: false},
	}

	// Parse the response body
	var obtained repository.User
	err = json.Unmarshal(rec.Body.Bytes(), &obtained)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", rec.Code)
	}

	// Compare the expected and obtained data
	if !reflect.DeepEqual(expected, obtained) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}

	// Using testify's assert package
	assert.Equal(t, expected, obtained, "The expected and obtained users should be equal")
	/*rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/user/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	// rctx := chi.NewRouteContext()
	// rctx.URLParams.Add("id", "1")
	// req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// handler := router.UserHandler{Repo: &mockUserModel{}}
	handler := handlers.UserHandler{Repo: &mockUserModel{}}
	http.HandlerFunc(handler.GetUserById).ServeHTTP(rec, req)
	expected := repository.User{
		ID:                 1,
		Name:               "Joe Doe",
		Email:              "johndoe@email.com",
		SubID:              "subid_123abc",
		VerificationStatus: true,
		SetupStatus:        "pending",
		CreatedAt:          pgtype.Timestamptz{Valid: false},
		UpdatedAt:          pgtype.Timestamptz{Valid: false},
	}

	var obtained repository.User
	err = json.Unmarshal(rec.Body.Bytes(), &obtained)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Compare the expected and obtained data
	if !reflect.DeepEqual(expected, obtained) {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}

	// Alternatively, using testify's assert
	assert.Equal(t, expected, obtained, "The expected and obtained users should be equal")
	// if expected != rec.Body.String() {
	// 		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	// 	}*/
}
