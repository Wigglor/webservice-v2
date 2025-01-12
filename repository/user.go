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

type Organization struct {
	ID                 int32              `json:"id"`
	Name               string             `json:"name"`
	SubscriptionId     string             `json:"subscriptionId"`
	PlanType           string             `json:"planType"`
	SubscriptionStatus string             `json:"subscriptionStatus"`
	NextBillingDate    pgtype.Timestamptz `json:"nextBillingDate"`
	CreatedAt          pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt          pgtype.Timestamptz `json:"updatedAt"`
}

type UserOrganization struct {
	UserId         string             `json:"userId"`
	OrganizationId string             `json:"organizationId"`
	Role           string             `json:"role"`
	CreatedAt      pgtype.Timestamptz `json:"createdAt"`
}

type ReturnOrgUser struct {
	Organization     Organization
	UserOrganization UserOrganization
}

type CreateOrganizationParams struct {
	UserId             int32              `json:"userId"`
	Name               string             `json:"name"`
	Role               string             `json:"role"`
	SubscriptionId     string             `json:"subscriptionId"`
	PlanType           string             `json:"planType"`
	SubscriptionStatus string             `json:"subscriptionStatus"`
	NextBillingDate    pgtype.Timestamptz `json:"nextBillingDate"`
}

type UserRepository interface {
	QueryAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	CheckUserBySubId(ctx context.Context, subId string) (User, error)
	QueryCreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	QueryCreateOrganization(ctx context.Context, arg CreateOrganizationParams) (ReturnOrgUser, error)
}
