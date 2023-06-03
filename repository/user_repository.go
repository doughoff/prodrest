package repository

import (
	"context"
	"fmt"

	"github.com/hoffax/prodrest/db"
	uuid "github.com/jackc/pgx-gofrs-uuid"
)

func (r *PgRepository) All(ctx context.Context) ([]*db.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			status,
			email,
			name,
			password,
			roles,
			created_at,
			updated_at
		FROM "users"
		ORDER BY 
		    created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*db.User, 0)
	for rows.Next() {
		user := db.User{}
		err := rows.Scan(
			&user.ID,
			&user.Status,
			&user.Email,
			&user.Name,
			&user.Password,
			&user.Roles,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("error while scanning user")
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PgRepository) GetByID(ctx context.Context, userID *uuid.UUID) (*db.User, error) {
	user := db.User{}
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			status,
			email,
			name,
			password,
			roles,
			created_at,
			updated_at
		FROM "users"
		WHERE
		    id = $1
	`, userID).Scan(
		&user.ID,
		&user.Status,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Roles,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PgRepository) NewUser(ctx context.Context, user *db.User) (*db.User, error) {
	var newUUID *uuid.UUID
	err := r.db.QueryRow(ctx, `
		insert into "users"(email, name, password, roles, created_at, updated_at) 
		values ($1, $2, $3, $4, now(), now()) 
		returning id
	`,
		user.Email,
		user.Name,
		user.Password,
		user.Roles,
	).Scan(&newUUID)
	if err != nil {
		return nil, err
	}

	newUser, err := r.GetByID(ctx, newUUID)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
