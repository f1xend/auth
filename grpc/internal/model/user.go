package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserInfo struct {
	Name     string
	Email    string
	Password string
	Role     bool
}

type UpdateUser struct {
	ID   int64
	Info UpdateUserInfo
}

type UpdateUserInfo struct {
	Name  string
	Email string
}
