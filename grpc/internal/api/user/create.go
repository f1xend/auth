package user

import (
	"context"
	"github.com/f1xend/auth/internal/converter"
	desc "github.com/f1xend/auth/pkg/auth_v1"
)

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.userService.Create(ctx, converter.ToUserInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
