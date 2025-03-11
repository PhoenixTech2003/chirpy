// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING id, created_at, updated_at, email, hashed_password
`

type CreateUserParams struct {
	Email          sql.NullString
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT  id, created_at, updated_at, email, hashed_password FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const updatePasswordAndEmail = `-- name: UpdatePasswordAndEmail :one
UPDATE users 
SET email = $1,
 hashed_password = $2,
 updated_at = NOW()
WHERE id = $3

RETURNING id, email, updated_at,created_at
`

type UpdatePasswordAndEmailParams struct {
	Email          sql.NullString
	HashedPassword string
	ID             uuid.UUID
}

type UpdatePasswordAndEmailRow struct {
	ID        uuid.UUID
	Email     sql.NullString
	UpdatedAt sql.NullTime
	CreatedAt sql.NullTime
}

func (q *Queries) UpdatePasswordAndEmail(ctx context.Context, arg UpdatePasswordAndEmailParams) (UpdatePasswordAndEmailRow, error) {
	row := q.db.QueryRowContext(ctx, updatePasswordAndEmail, arg.Email, arg.HashedPassword, arg.ID)
	var i UpdatePasswordAndEmailRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
