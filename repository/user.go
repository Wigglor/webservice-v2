package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID                 int32              `json:"id"`
	Name               string             `json:"name"`
	Email              string             `json:"email"`
	SubID              string             `json:"subId"`
	VerificationStatus bool               `json:"verificationStatus"`
	SetupStatus        string             `json:"setupStatus"`
	CreatedAt          pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt          pgtype.Timestamptz `json:"updatedAt"`
}

type CreateUserParams struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	SubId              string `json:"subId"`
	VerificationStatus bool   `json:"verificationStatus"`
	SetupStatus        string `json:"setupStatus"`
}

type UserRepository interface {
	QueryAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	CheckUserBySubId(ctx context.Context, subId string) (User, error)
	QueryCreateUser(ctx context.Context, arg CreateUserParams) (User, error)
}
