package repository

import "github.com/jackc/pgx/v5"

type PgRepository struct {
	db *pgx.Conn
}

func NewPgRepository(conn *pgx.Conn) *PgRepository {
	return &PgRepository{
		db: conn,
	}
}
