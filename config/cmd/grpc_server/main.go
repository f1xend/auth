package main

import (
	"context"
	"database/sql"
	"flag"
	sq "github.com/Masterminds/squirrel"
	"github.com/f1xend/auth/config/internal/config"
	"github.com/f1xend/auth/config/internal/config/env"
	desc "github.com/f1xend/auth/config/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// Делаем запрос на вставку записи в таблицу auth
	builderInsert := sq.Insert("user").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email, password, role").
		Values(
			req.Info.Name,
			req.Info.Email,
			req.Info.Password,
			req.Info.Role,
		)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	if err = s.pool.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id %d", userID)

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	builderSelectOne := sq.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
		From("user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	var name, email, password string
	var createdAt time.Time
	var updatedAt sql.NullTime
	var role bool

	if err = s.pool.QueryRow(ctx, query, args...).
		Scan(&userID, &name, &email, &password, &role, &createdAt, &updatedAt); err != nil {
		log.Fatalf("failed to select user %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s, password: %s, created_at: %s, updated_at: %s",
		userID, name, email, password, createdAt, updatedAt)

	var updateAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updateAtTime = timestamppb.New(updatedAt.Time)
	}

	var roleAdmin desc.Role
	if role == true {
		roleAdmin = desc.Role_admin
	} else {
		roleAdmin = desc.Role_user
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id: userID,
			Info: &desc.UserInfo{
				Name:            name,
				Email:           email,
				Password:        password,
				PasswordConfirm: password,
				Role:            roleAdmin,
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updateAtTime,
		},
	}, nil
}

//func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error)

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

	// создаю сервер
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
