-- +goose Up
CREATE TABLE chirps(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid,
    body TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose Down
DROP TABLE chirps;