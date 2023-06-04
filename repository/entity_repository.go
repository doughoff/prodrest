package repository

import (
	"context"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"time"
)

type Entity struct {
	ID        *pgxuuid.UUID
	Status    string
	Name      string
	RUC       string
	CI        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FetchEntitiesParams struct {
	StatusOptions []string
	Search        string
	Limit         int
	Offset        int
}

type FetchEntitiesResult struct {
	TotalCount int
	Items      []*Entity
}

func (r *PgRepository) FetchEntities(ctx context.Context, param *FetchEntitiesParams) (*FetchEntitiesResult, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
		    COUNT(*) OVER() AS full_count,
			id,
			status,
			name,
			ruc,
			ci,
			created_at,
			updated_at
		FROM "entities"
		WHERE
		    status = ANY($1::status[])
			 OR (
			     name ILIKE '%' || $2 || '%'
			     OR ruc ILIKE '%' || $2 || '%'
				 OR ci ILIKE '%' || $2 || '%'
			 )
		ORDER BY
		    created_at DESC
		LIMIT $3
		OFFSET $4
	`, param.StatusOptions, param.Search, param.Limit, param.Offset)
	if err != nil {
		return nil, err
	}

	var result FetchEntitiesResult
	for rows.Next() {
		var item Entity
		err := rows.Scan(
			&result.TotalCount,
			&item.ID,
			&item.Status,
			&item.Name,
			&item.RUC,
			&item.CI,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Items = append(result.Items, &item)
	}

	return &result, nil
}

func (r *PgRepository) GetEntityById(ctx context.Context, id *pgxuuid.UUID) (*Entity, error) {
	var item Entity
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			status,
			name,
			ruc,
			ci,
			created_at,
			updated_at
		FROM "entities"
		WHERE id = $1
	`, id).Scan(
		&item.ID,
		&item.Status,
		&item.Name,
		&item.RUC,
		&item.CI,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

type CreateEntityParams struct {
	Name string
	RUC  string
	CI   string
}

func (r *PgRepository) CreateEntity(ctx context.Context, param *CreateEntityParams) (*Entity, error) {
	var item Entity
	err := r.db.QueryRow(ctx, `
		INSERT INTO "entities" (
			name,
			ruc,
			ci
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING
			id,
			status,
			name,
			ruc,
			ci,
			created_at,
			updated_at
	`, param.Name, param.RUC, param.CI).Scan(
		&item.ID,
		&item.Status,
		&item.Name,
		&item.RUC,
		&item.CI,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

type UpdateEntityParams struct {
	ID     *pgxuuid.UUID
	Status string
	Name   string
	RUC    string
	CI     string
}

func (r *PgRepository) UpdateEntity(ctx context.Context, param *UpdateEntityParams) (*Entity, error) {
	var entity Entity
	err := r.db.QueryRow(ctx, `
		UPDATE "entities" SET
			status = $2,
			name = $3,
			ruc = $4,
			ci = $5
		WHERE id = $1
		RETURNING
			id,
			status,
			name,
			ruc,
			ci,
			created_at,
			updated_at
	`, param.ID, param.Status, param.Name, param.RUC, param.CI).Scan(
		&entity.ID,
		&entity.Status,
		&entity.Name,
		&entity.RUC,
		&entity.CI,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *PgRepository) GetEntityByRUC(ctx context.Context, ruc string) (*Entity, error) {
	var item Entity
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			status,
			name,
			ruc,
			ci,
			created_at,
			updated_at
		FROM "entities"
		WHERE ruc = $1
	`, ruc).Scan(
		&item.ID,
		&item.Status,
		&item.Name,
		&item.RUC,
		&item.CI,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *PgRepository) GetEntityByCI(ctx context.Context, ci string) (*Entity, error) {
	var item Entity
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			status,
			name,
			ruc,
			ci,
			created_at,
			updated_at
		FROM "entities"
		WHERE ci = $1
	`, ci).Scan(
		&item.ID,
		&item.Status,
		&item.Name,
		&item.RUC,
		&item.CI,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
