package user

import (
	"context"
	"github.com/f1xend/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	id, err := s.userRepository.Create(ctx, info)
	if err != nil {
		return 0, err
	}
	return id, nil
}
