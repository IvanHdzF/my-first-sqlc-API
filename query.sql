-- Example queries for sqlc
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	username VARCHAR(30) NOT NULL,
	bio VARCHAR(400),
	avatar VARCHAR(200),
	phone VARCHAR(25),
	email VARCHAR(40),
	password VARCHAR(50),
	status VARCHAR(15),
	CHECK(COALESCE(phone, email) IS NOT NULL)
);

-- name: GetAuthor :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM users
ORDER BY id;

-- name: CreateAuthor :one
INSERT INTO users (
  username, bio, avatar,phone,email,password,status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM users
WHERE id = $1;