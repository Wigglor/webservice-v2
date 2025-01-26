package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// func (m UserRepo) QueryAllUsers() ([]User, error) {
func (m *UserRepo) QueryAllUsers(ctx context.Context) ([]User, error) {
	rows, err := m.db.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.SubID,
			&i.VerificationStatus,
			&i.SetupStatus,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// func (m UserRepo) GetUserByID(id int32) (User, error) {
func (m *UserRepo) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := m.db.QueryRow(ctx, `-- name: GetUserByID :one
SELECT
  id,
  name,
  email,
  sub_id,
  verification_status,
  setup_status,
  created_at,
  updated_at
FROM
  users
WHERE
  id = $1
`, id)
	// 	row := m.db.QueryRow(ctx, `-- name: GetUserByID :one
	// SELECT
	//   id,
	//   name,
	//   email,
	//   sub_id,
	//   verification_status,
	//   setup_status,
	//   created_at,
	//   updated_at
	// FROM
	//   user_tests
	// WHERE
	//   id = $1
	// `, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.SubID,
		&i.VerificationStatus,
		&i.SetupStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func (m *UserRepo) CheckUserBySubId(ctx context.Context, subId string) (User, error) {
	row := m.db.QueryRow(ctx, `-- name: CheckUserBySubId :one
SELECT
  id,
  name,
  email,
  sub_id,
  verification_status,
  setup_status,
  created_at,
  updated_at
FROM
  users
WHERE
  sub_id = $1
`, subId)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.SubID,
		&i.VerificationStatus,
		&i.SetupStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// func (m UserRepo) QueryCreateUser(arg CreateUserParams) (User, error) {
func (m *UserRepo) QueryCreateUser(ctx context.Context, arg CreateUserParams, subId string) (User, error) {
	row := m.db.QueryRow(ctx, `-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  sub_id,
  verification_status,
  setup_status
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING
  id,
  name,
  email,
  sub_id,
  verification_status,
  setup_status,
  created_at,
  updated_at
`,
		arg.Name,
		arg.Email,
		// arg.SubId,
		subId,
		arg.VerificationStatus,
		// arg.SetupStatus,
		"in_progress", // this should always be in_progress at first since user has not yet set up organizations and user_organizations
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.SubID,
		&i.VerificationStatus,
		&i.SetupStatus,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func (m *UserRepo) QueryCreateOrganization(ctx context.Context, arg CreateOrganizationParams) (ReturnOrgUser, error) {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return ReturnOrgUser{}, fmt.Errorf("begin transaction: %w", err)
	}

	// If anything goes wrong, roll back
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var org Organization
	insertOrgQuery := `
        INSERT INTO organizations (name, subscription_id, plan_type, subscription_status, next_billing_date)
        VALUES ($1, $2, $3, $4, $5)
		RETURNING
        id, 
		name, 
		subscription_id, 
		plan_type, 
		subscription_status, 
		next_billing_date,
		created_at,
  		updated_at
    `
	err = tx.QueryRow(ctx, insertOrgQuery, arg.Name, arg.SubscriptionId, arg.PlanType, arg.SubscriptionStatus, arg.NextBillingDate).Scan(&org.ID, &org.Name, &org.SubscriptionId, &org.PlanType, &org.SubscriptionStatus, &org.NextBillingDate, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return ReturnOrgUser{}, fmt.Errorf("insert organization: %w", err)
	}

	// 2) Insert the membership into user_organizations
	insertMembershipQuery := `
	 INSERT INTO user_organizations (user_id, organization_id, role)
	 VALUES ($1, $2, $3)
	 RETURNING user_id, organization_id, role, created_at
 `
	var userOrg UserOrganization
	err = tx.QueryRow(ctx, insertMembershipQuery, arg.UserId, org.ID, arg.Role).Scan(&userOrg.UserId, &userOrg.OrganizationId, &userOrg.Role, &userOrg.CreatedAt)
	if err != nil {
		return ReturnOrgUser{}, fmt.Errorf("insert user_organization link: %w", err)
	}

	// 3) Update user’s setup_status to "completed"
	updateUserSetupQuery := `
		UPDATE users
		SET setup_status = 'completed'
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, updateUserSetupQuery, arg.UserId)
	if err != nil {
		return ReturnOrgUser{}, fmt.Errorf("update user setup_status: %w", err)
	}

	// 3) Commit if everything succeeded
	if err = tx.Commit(ctx); err != nil {
		return ReturnOrgUser{}, fmt.Errorf("commit transaction: %w", err)
	}

	return ReturnOrgUser{
		org,
		userOrg,
	}, err
}

// type PostgresUserRepository struct {
// 	DB *sql.DB
// }

// // GetByID fetches a user by their ID
// func (r *PostgresUserRepository) GetByID(id int) (*models.User, error) {
// 	user := &models.User{}
// 	query := "SELECT id, name, email FROM users WHERE id = $1"
// 	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// // Create inserts a new user into the database
// func (r *PostgresUserRepository) Create(user *models.User) error {
// 	query := "INSERT INTO users (name, email) VALUES ($1, $2)"
// 	_, err := r.DB.Exec(query, user.Name, user.Email)
// 	return err
// }

// // Update modifies an existing user's information
// func (r *PostgresUserRepository) Update(user *models.User) error {
// 	query := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
// 	_, err := r.DB.Exec(query, user.Name, user.Email, user.ID)
// 	return err
// }

// // Delete removes a user by their ID
// func (r *PostgresUserRepository) Delete(id int) error {
// 	query := "DELETE FROM users WHERE id = $1"
// 	_, err := r.DB.Exec(query, id)
// 	return err
// }
