package user

import (
	"github.com/f1xend/auth/internal/service"
	desc "github.com/f1xend/auth/pkg/auth_v1"
)

type Server struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}
