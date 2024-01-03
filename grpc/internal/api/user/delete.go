package user

import (
	"context"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := s.userService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
