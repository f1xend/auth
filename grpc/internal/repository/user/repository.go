package user

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/f1xend/auth/internal/repository"
	"github.com/f1xend/auth/internal/repository/user/converter"
	"github.com/f1xend/auth/internal/repository/user/model"
	desc "github.com/f1xend/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, req *desc.UserInfo) (int64, error) {
	// Делаем запрос на вставку записи в таблицу auth
	var role bool
	if req.Role == desc.Role_admin {
		role = true
	}
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(
			req.Name,
			req.Email,
			req.Password,
			role,
		).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userID int64
	if err = r.db.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id %d", userID)

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.User, error) {
	log.Printf("User id: %d", id)

	builderSelectOne := sq.Select(
		idColumn, nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	info := model.Info{}
	user := model.User{
		Info: &info,
	}

	log.Println(ctx, query, args)

	err = r.db.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.Password, &user.Info.Role,
			&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Fatalf("failed to select user %v", err)
	}
	//
	//log.Printf("id: %d, name: %s, email: %s, password: %s, created_at: %s, updated_at: %s",
	//	user.ID, user.Info.Name, user.Info.Email, user.Info.Password,
	//	user.CreatedAt, user.UpdatedAt)

	return converter.ToUserFromRepo(&user), nil
}
