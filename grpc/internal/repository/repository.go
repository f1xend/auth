package repository

import (
	"context"
	"github.com/f1xend/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, req *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}
