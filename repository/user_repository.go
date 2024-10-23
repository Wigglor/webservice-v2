package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateUserParams struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	SubID              string `json:"subId"`
	VerificationStatus bool   `json:"verificationStatus"`
	SetupStatus        string `json:"setupStatus"`
}

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
  user_tests
WHERE
  id = $1
`, id)
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
func (m *UserRepo) QueryCreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := m.db.QueryRow(ctx, `-- name: CreateUser :one
INSERT INTO user_tests (
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
		arg.SubID,
		arg.VerificationStatus,
		arg.SetupStatus,
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
