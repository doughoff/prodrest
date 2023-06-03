package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"os"
	"time"

	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
)

type MyTracer struct {
	logger *log.Logger
}

func (t *MyTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	fmt.Printf("sql: %v\n args:%v\n", data.SQL, data.Args)
	return nil
}
func (t *MyTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	fmt.Printf("err:%v\n tag:%v\n", data.Err, data.CommandTag)
	return
}

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatalf("Error closing database connection")
		}
	}(conn, context.Background())

	connConfig, err := pgx.ParseConfig(os.Getenv("DB_URL"))
	connConfig.AfterConnect = func(ctx context.Context, pgconn *pgconn.PgConn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	connConfig.Tracer = &MyTracer{}

	repo := repository.NewPgRepository(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	allStatus := []string{"ACTIVE"}
	users, err := repo.GetAllUsers(ctx, &repository.GetAllUsersParams{
		Status: allStatus,
	})
	if err != nil {
		log.Fatalf("Unable to get users: %v\n", err)
	}

	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}

}
