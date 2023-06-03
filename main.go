package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
)

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	pgxuuid.Register(conn.TypeMap())

	repo := repository.NewPgRepository(conn)

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Microsecond)

	users, err := repo.All(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get users: %v\n", err)
		os.Exit(1)
	}

	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}

}
