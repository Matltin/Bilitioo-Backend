// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (
  "email",
  "phone_number",
  "hashed_password",
  "email_verified",
  "phone_verified"
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, email, phone_number, hashed_password, password_change_at, role, status, phone_verified, email_verified, created_at
`

type CreateUserParams struct {
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	HashedPassword string `json:"hashed_password"`
	EmailVerified  bool   `json:"email_verified"`
	PhoneVerified  bool   `json:"phone_verified"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.PhoneNumber,
		arg.HashedPassword,
		arg.EmailVerified,
		arg.PhoneVerified,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PhoneNumber,
		&i.HashedPassword,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Status,
		&i.PhoneVerified,
		&i.EmailVerified,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT 
    "id", "email", "phone_number", "hashed_password", "email_verified", "phone_verified"
FROM "user"
WHERE "email" = $1 OR "phone_number" = $2
`

type GetUserParams struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type GetUserRow struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	HashedPassword string `json:"hashed_password"`
	EmailVerified  bool   `json:"email_verified"`
	PhoneVerified  bool   `json:"phone_verified"`
}

func (q *Queries) GetUser(ctx context.Context, arg GetUserParams) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, arg.Email, arg.PhoneNumber)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PhoneNumber,
		&i.HashedPassword,
		&i.EmailVerified,
		&i.PhoneVerified,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT 
    "id", "email", "phone_number", "hashed_password"
FROM "user"
WHERE "id" = $1
`

type GetUserByIDRow struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PhoneNumber,
		&i.HashedPassword,
	)
	return i, err
}

const initialProfile = `-- name: InitialProfile :exec
INSERT INTO "profile" (
  "user_id"
) VALUES (
  $1
)
`

func (q *Queries) InitialProfile(ctx context.Context, userID int64) error {
	_, err := q.db.ExecContext(ctx, initialProfile, userID)
	return err
}

const updateUserContact = `-- name: UpdateUserContact :one
UPDATE "user"
SET
  email = COALESCE($2, email),
  phone_number = COALESCE($3, phone_number),
  hashed_password = COALESCE($4, hashed_password)
WHERE id = $1
RETURNING id, email, phone_number, hashed_password, password_change_at, role, status, phone_verified, email_verified, created_at
`

type UpdateUserContactParams struct {
	ID             int64          `json:"id"`
	Email          sql.NullString `json:"email"`
	PhoneNumber    sql.NullString `json:"phone_number"`
	HashedPassword sql.NullString `json:"hashed_password"`
}

func (q *Queries) UpdateUserContact(ctx context.Context, arg UpdateUserContactParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserContact,
		arg.ID,
		arg.Email,
		arg.PhoneNumber,
		arg.HashedPassword,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PhoneNumber,
		&i.HashedPassword,
		&i.PasswordChangeAt,
		&i.Role,
		&i.Status,
		&i.PhoneVerified,
		&i.EmailVerified,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserEmailVerified = `-- name: UpdateUserEmailVerified :exec
UPDATE "user"
SET email_verified = $2
WHERE id = $1
`

type UpdateUserEmailVerifiedParams struct {
	ID            int64 `json:"id"`
	EmailVerified bool  `json:"email_verified"`
}

func (q *Queries) UpdateUserEmailVerified(ctx context.Context, arg UpdateUserEmailVerifiedParams) error {
	_, err := q.db.ExecContext(ctx, updateUserEmailVerified, arg.ID, arg.EmailVerified)
	return err
}
