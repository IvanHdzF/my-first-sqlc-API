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

-- name: GetUser :one
select jsonb_build_object(
	'id',id, 
	'username',username, 
	'avatar',avatar,
	'phone',phone,
	'email',email,
	'password',password,
	'status',status)
from (
	SELECT * FROM USERS WHERE ID=$1
)AS selectedUser;

-- name: ListUsers :one
select jsonb_agg(jsonb_build_object(
	'id',id, 
	'username',username, 
	'avatar',avatar,
	'phone',phone,
	'email',email,
	'password',password,
	'status',status))
from (
	SELECT * FROM USERS ORDER BY id ASC
)AS sortedUser;

-- name: CreateUser :one
INSERT INTO users (
  username, bio, avatar,phone,email,password,status
) 
SELECT 
	username,
	bio,
	avatar,
	phone,
	email,
	password,
	status
FROM jsonb_populate_record(null::users, sqlc.arg(payload)) RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE USERS
SET 
	username = $2,
	bio=$3,
	avatar=$4,
	phone=$5,
	email=$6,
	password=$7,
	status=$8
WHERE ID=$1;
