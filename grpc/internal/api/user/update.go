package user

import (
	"context"
	"github.com/f1xend/auth/internal/converter"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	err := s.userService.Update(ctx, converter.ToUserUpdateFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
