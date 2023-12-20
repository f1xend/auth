package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

const grpcPort = 50051

type Server struct {
	desc.UnimplementedUserV1Server
}

func (s *Server) Get(context context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	pass := gofakeit.Password(true, true, true, true, false, 6)
	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:            gofakeit.Name(),
				Email:           gofakeit.Email(),
				Password:        pass,
				PasswordConfirm: pass,
				Role:            desc.Role_admin,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
