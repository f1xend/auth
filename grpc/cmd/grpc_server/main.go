package main

import (
	"context"
	"flag"
	"github.com/f1xend/auth/internal/config"
	"github.com/f1xend/auth/internal/config/env"
	"github.com/f1xend/auth/internal/repository"
	"github.com/f1xend/auth/internal/repository/user"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "prod.env", "path to config file")
}

type Server struct {
	desc.UnimplementedUserV1Server
	userRepository repository.UserRepository
}

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.userRepository.Create(ctx, req.GetInfo())
	if err != nil {

		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: userObj,
	}, nil
}

func main() {
	flag.Parse()

	ctx := context.Background()

	// Считываю переменные окружения
	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	// создаю пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	// создаю репозиторий
	userRepo := user.NewRepository(pool)

	// создаю сервер
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &Server{userRepository: userRepo})

	log.Printf("Server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
