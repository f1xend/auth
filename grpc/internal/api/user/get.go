package user

import (
	"context"
	"github.com/f1xend/auth/internal/converter"
	desc "github.com/f1xend/auth/pkg/auth_v1"
)

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := s.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(userObj),
	}, nil
}
