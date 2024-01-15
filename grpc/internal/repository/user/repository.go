package user

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/f1xend/auth/internal/model"
	"github.com/f1xend/auth/internal/repository"
	"github.com/f1xend/auth/internal/repository/user/converter"
	modelRepo "github.com/f1xend/auth/internal/repository/user/model"
	"github.com/f1xend/platform-common/pkg/db"
	"github.com/pkg/errors"
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

func (r *repo) Update(ctx context.Context, req *model.UpdateUser) error {
	log.Printf("update user id: %d, userInfo: %v", req.ID, req.Info)

	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, req.Info.Name).
		Set(emailColumn, req.Info.Email).
		Where(sq.Eq{idColumn: req.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	log.Printf("delete user id: %d", id)

	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build delete query")
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}
