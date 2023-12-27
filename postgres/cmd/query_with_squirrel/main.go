package main

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database %v", err)
	}
	defer pool.Close()

	//билд инсерта
	buildInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password").
		Values(gofakeit.Name(),
			gofakeit.Email(),
			gofakeit.Password(true, true, true, true, false, 10),
		).
		Suffix("RETURNING id")

	query, args, err := buildInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query %v", err)
	}

	var userID int
	err = pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)

	builderSelect := sq.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select users")
	}

	var id int
	var name, email, password string
	var createdAt time.Time
	var updatedAt sql.NullTime
	var role bool

	for rows.Next() {
		if err = rows.Scan(&id, &name, &email, &password, &role, &createdAt, &updatedAt); err != nil {
			log.Fatalf("failed to scan user: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, password: %s, created_at: %s, updated_at: %s",
			id, name, email, password, createdAt, updatedAt)
	}

	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("updated %v rows", res.RowsAffected())

	builderSelectOne := sq.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatal("failed to build query: %v", err)
	}

	if err = pool.QueryRow(ctx, query, args...).
		Scan(&id, &name, &email, &password, &role, &createdAt, &updatedAt); err != nil {
		log.Fatalf("failed to select user: %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s, password: %s, role: %s, created_at: %s, updated_at: %s",
		id, name, email, password, role, createdAt, updatedAt)
}
