// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
)

type User struct {
	ID        int32
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	Username  string
	Bio       sql.NullString
	Avatar    sql.NullString
	Phone     sql.NullString
	Email     sql.NullString
	Password  sql.NullString
	Status    sql.NullString
}
