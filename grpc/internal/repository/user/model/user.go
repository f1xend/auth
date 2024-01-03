package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `db:"id"`
	Info      UserInfo     `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type UserInfo struct {
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     bool   `db:"role"`
}
