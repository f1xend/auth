package service

import (
	"context"
	"github.com/f1xend/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, req *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, req *model.UpdateUser) error
	Delete(ctx context.Context, id int64) error
}
