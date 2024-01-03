package user

import (
	"context"
	"github.com/f1xend/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	//log.Printf("id: %d, name: %s, email: %s, password: %s, created_at: %s, updated_at: %s",
	//	user.ID, user.Info.Name, user.Info.Email, user.Info.Password,
	//	user.CreatedAt, user.UpdatedAt)
	return user, nil
}
