package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/storage/memory"
	"github.com/hoffax/prodrest/middleware"
	"github.com/hoffax/prodrest/repository"
	"github.com/hoffax/prodrest/routes"
	"github.com/hoffax/prodrest/services"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"os"
	"time"
)

func Serve() {
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

	repo := repository.NewPgRepository(conn)
	sm, err := services.NewServiceManager(repo)
	if err != nil {
		log.Fatalf("Could not open service manager\n %v", err)
	}

	memoryStore := memory.New(memory.Config{
		GCInterval: 5 * time.Second,
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.FiberCustomErrorHandler,
	})
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(middleware.AuthMiddleware(memoryStore))
	app.Use(cache.New(cache.Config{
		Expiration:   1 * time.Second,
		CacheControl: true,
	}))

	handlers := routes.NewHandlers(app, sm, memoryStore)
	app.Get("/metrics", monitor.New())
	handlers.RegisterAuthRoutes()
	handlers.RegisterUserRoutes()

	err = app.Listen(":3088")
	if err != nil {
		log.Fatalf("err listen: %v", err)
	}

	return
}
