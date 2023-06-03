package repository

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/jackc/pgx-gofrs-uuid"
)

type User struct {
	ID        *uuid.UUID
	Status    string
	Roles     []string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetAllUsersParams struct {
	Status []string
}

func (r *PgRepository) GetAllUsers(ctx context.Context, params *GetAllUsersParams) ([]*User, error) {
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
		WHERE 
		    status = ANY($1::status[])
		ORDER BY 
		    created_at DESC
	`, params.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := User{}
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

func (r *PgRepository) GetUserByID(ctx context.Context, userID *uuid.UUID) (*User, error) {
	user := User{}
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

func (r *PgRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}
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
	`, email).Scan(
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

type NewUserParams struct {
	Email    string
	Name     string
	Password string
	Roles    []string
}

func (r *PgRepository) NewUser(ctx context.Context, user *NewUserParams) (*User, error) {
	newUser := &User{}
	err := r.db.QueryRow(ctx, `
		insert into "users"(email, name, password, roles) 
		values ($1, $2, $3, $4) 
		returning id, status, email, name, password, roles, created_at, updated_at
	`,
		user.Email,
		user.Name,
		user.Password,
		user.Roles,
	).Scan(
		&newUser.ID,
		&newUser.Status,
		&newUser.Email,
		&newUser.Name,
		&newUser.Password,
		&newUser.Roles,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

type UpdateUserParams struct {
	ID       *uuid.UUID
	status   []string
	Email    string
	Name     string
	Password string
	roles    []string
}

func (r *PgRepository) UpdateUser(ctx context.Context, params *UpdateUserParams) (*User, error) {
	updatedUser := &User{}
	err := r.db.QueryRow(ctx, `
		update "users"
		set 
			status = $2,
			email = $3,
			name = $4,
			password = $5,
			roles = $6,
			updated_at = now()
		where 
		    id = $1
		returning id, status, email, name, password, roles, created_at, updated_at
	`,
		params.ID,
		params.status,
		params.Email,
		params.Name,
		params.Password,
		params.roles,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Status,
		&updatedUser.Email,
		&updatedUser.Name,
		&updatedUser.Password,
		&updatedUser.Roles,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
