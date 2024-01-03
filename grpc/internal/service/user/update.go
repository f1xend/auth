package user

import (
	"context"
	"github.com/f1xend/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, user *model.UpdateUser) error {
	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
