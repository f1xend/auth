package user

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/f1xend/auth/internal/client/db"
	"github.com/f1xend/auth/internal/model"
	"github.com/f1xend/auth/internal/repository"
	"github.com/f1xend/auth/internal/repository/user/converter"
	modelRepo "github.com/f1xend/auth/internal/repository/user/model"
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
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, req *model.UserInfo) (int64, error) {
	// Делаем запрос на вставку записи в таблицу auth
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(
			req.Name,
			req.Email,
			req.Password,
			req.Role,
		).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var userID int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID); err != nil {
		return 0, err
	}

	log.Printf("inserted user with id %d", userID)

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	log.Printf("User id: %d", id)

	builderSelectOne := sq.Select(
		idColumn, nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	if err = r.db.DB().ScanOneContext(ctx, &user, q, args...); err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}
