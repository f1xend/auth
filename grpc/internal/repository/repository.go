package repository

import (
	"context"
	desc "github.com/f1xend/auth/pkg/auth_v1"
)

type UserRepository interface {
	Create(ctx context.Context, req *desc.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*desc.User, error)
}
