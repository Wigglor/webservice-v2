package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Wigglor/webservice-v2/repository"
	"github.com/Wigglor/webservice-v2/router"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

type mockUserModel struct{}

func (m *mockUserModel) QueryAllUsers(ctx context.Context) ([]repository.User, error) {
	var users []repository.User

	users = append(users, repository.User{ID: 1, Name: "Joe Doe", Email: "johndoe@email.com", SubID: "subid_123abc", VerificationStatus: true, SetupStatus: "pending"})
	users = append(users, repository.User{ID: 2, Name: "Jane Doe", Email: "janedoe@email.com", SubID: "subid_456def", VerificationStatus: true, SetupStatus: "pending"})

	return users, nil
}

func (m *mockUserModel) GetUserByID(ctx context.Context, id int32) (repository.User, error) {
	// var users []repository.User

	users := repository.User{ID: 1, Name: "Joe Doe", Email: "johndoe@email.com", SubID: "subid_123abc", VerificationStatus: true, SetupStatus: "pending"}
	// users = append(users, repository.User{ID: 2, Name: "Jane Doe", Email: "janedoe@email.com", SubID: "subid_456def", VerificationStatus: true, SetupStatus: "pending"})

	return users, nil
}
func TestGetUsers(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	handler := router.UserHandler{Repo: &mockUserModel{}}
	http.HandlerFunc(handler.GetUsers).ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())
	expected := []repository.User{
		{
			ID:                 1,
			Name:               "Joe Doe",
			Email:              "johndoe@email.com",
			SubID:              "subid_123abc",
			VerificationStatus: true,
			SetupStatus:        "pending",
			CreatedAt:          pgtype.Timestamptz{Valid: false},
			UpdatedAt:          pgtype.Timestamptz{Valid: false},
		},
		{
			ID:                 2,
			Name:               "Jane Doe",
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
}

/*
func TestGetUsers(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	userRepo := repository.NewUserRepository(&mockUserModel{})
	// env := Env{users: &mockUserModel{}, context: context.Background()}
	env := repository.UserRepo{db: &mockUserModel{}}

	http.HandlerFunc(env.GetUsers).ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())
	expected := []repository.User{
		{
			ID:                 1,
			Name:               "Joe Doe",
			Email:              "johndoe@email.com",
			SubID:              "subid_123abc",
			VerificationStatus: true,
			SetupStatus:        "pending",
			CreatedAt:          pgtype.Timestamptz{Valid: false},
			UpdatedAt:          pgtype.Timestamptz{Valid: false},
		},
		{
			ID:                 2,
			Name:               "Jane Doe",
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
	// 	t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	// }

}*/
