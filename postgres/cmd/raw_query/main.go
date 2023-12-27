package main

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	pgx "github.com/jackc/pgx/v4"
	"log"
	"time"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database %v", err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}()

	// insert to table
	res, err := conn.Exec(ctx, "INSERT INTO users (name, email, password) values ($1, $2, $3)",
		gofakeit.Name(),
		gofakeit.Email(),
		gofakeit.Password(true, true, true, true, false, 10),
	)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("inserted %v rows", res.RowsAffected())

	//Запрос на выборку всех записей из таблицы users
	rows, err := conn.Query(ctx, "SELECT id, name, email, password, role, created_at, updated_at from users")
	if err != nil {
		log.Fatalf("failed to select users %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email, password string
		var role bool
		var createdAt time.Time
		var UpdatedAt sql.NullTime

		if err = rows.Scan(&id, &name, &email, &password, &role, &createdAt, &UpdatedAt); err != nil {
			log.Fatalf("failed to scan users %v", err)
		}

		log.Printf("id: %+d name: %+s email: %+s password: %+s role: %+t createdAt: %+v UpdatedAt: %+v",
			id, name, email, password, role, createdAt, UpdatedAt)
	}
}
