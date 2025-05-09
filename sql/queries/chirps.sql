-- name: CreateChirp :one
INSERT INTO chirps (id, user_id, body) 
VALUES
(
    gen_random_uuid(),
    $1,
    $2
)


RETURNING *;

-- name: DeleteUserChirp :exec
DELETE FROM chirps WHERE id = $1 AND user_id = $2;


-- name: GetAllChirpsByAuthor :many
SELECT * FROM chirps WHERE user_id = $1;

-- name: GetAllChirps :many
SELECT * FROM chirps;
