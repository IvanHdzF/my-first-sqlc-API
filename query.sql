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

CREATE TABLE posts (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    url character varying(200) NOT NULL,
    caption character varying(240),
    lat real,
    lng real,
    user_id integer NOT NULL,
    CONSTRAINT posts_lat_check CHECK (((lat IS NULL) OR ((lat >= ('-90'::integer)::double precision) AND (lat <= (90)::double precision)))),
    CONSTRAINT posts_lng_check CHECK (((lng IS NULL) OR ((lng >= ('-180'::integer)::double precision) AND (lng <= (180)::double precision))))
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
FROM jsonb_populate_record(null::users, @payload) RETURNING id;

-- name: DeleteUser :one
DELETE FROM users
WHERE users.id=(SELECT id FROM jsonb_populate_record(null::users, sqlc.arg(payload)) AS jsonRequest) 
RETURNING id;
-- name: UpdateUser :exec
UPDATE USERS
SET (username,bio,avatar,phone,email,password,status)= (SELECT username,bio,avatar,phone,email,password,status 
														FROM jsonb_populate_record(null::users, $1))
WHERE users.id=sqlc.arg(id);

-- name: GetUserPosts :many
SELECT username, url,caption
FROM posts AS p
JOIN users AS u ON p.user_id=u.id
WHERE p.id=(SELECT id FROM jsonb_populate_record(null::users, @payload));