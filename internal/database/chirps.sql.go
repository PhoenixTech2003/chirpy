// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: chirps.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps (id, user_id, body) 
VALUES
(
    gen_random_uuid(),
    $1,
    $2
)

RETURNING id, user_id, body
`

type CreateChirpParams struct {
	UserID uuid.NullUUID
	Body   sql.NullString
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.UserID, arg.Body)
	var i Chirp
	err := row.Scan(&i.ID, &i.UserID, &i.Body)
	return i, err
}
