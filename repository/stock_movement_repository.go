package repository

import (
	"context"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"time"
)

type StockMovement struct {
	ID        *pgxuuid.UUID
	Status    string
	Type      string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	EntityID       *pgxuuid.UUID
	EntityName     string
	EntityDocument string

	CreatedByUserID     *pgxuuid.UUID
	CreatedByUserName   string
	CancelledByUserID   *pgxuuid.UUID
	CancelledByUserName string
}

type FetchStockMovementsParams struct {
	StatusOptions []string
	TypeOptions   []string
	StartDate     time.Time
	Limit         int
	Offset        int
}

type FetchStockMovementsResult struct {
	TotalCount int
	Items      []*StockMovement
}

func (r *PgRepository) FetchStockMovements(ctx context.Context, param *FetchStockMovementsParams) (*FetchStockMovementsResult, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
		    COUNT(*) OVER() AS full_count,
			sm.id,
			sm.status,
			sm.type,
			sm.date,
			sm.created_at,
			sm.updated_at,
			e.id,
			coalesce(e.name, '-'),
			coalesce(e.ruc, e.ci, '-'),
			cu.id,
			cu.name,
			cu2.id,
			coalesce(cu2.name, '-')
		FROM "stock_movements" sm
		LEFT JOIN "entities" e ON e.id = sm.entity_id
		LEFT JOIN "users" cu ON cu.id = sm.created_by_user_id
		LEFT JOIN "users" cu2 ON cu2.id = sm.cancelled_by_user_id
		WHERE
		    sm.status = ANY($1::status[])
			AND sm.type = ANY($2::movement_type[])
			AND sm.date >= $3
		ORDER BY
		    sm.created_at DESC
		LIMIT $4
		OFFSET $5
	`, param.StatusOptions, param.TypeOptions, param.StartDate, param.Limit, param.Offset)
	if err != nil {
		return nil, err
	}

	var result FetchStockMovementsResult
	for rows.Next() {
		sm := StockMovement{}
		rows.FieldDescriptions()
		err := rows.Scan(
			&result.TotalCount,
			&sm.ID,
			&sm.Status,
			&sm.Type,
			&sm.Date,
			&sm.CreatedAt,
			&sm.UpdatedAt,
			&sm.EntityID,
			&sm.EntityName,
			&sm.EntityDocument,
			&sm.CreatedByUserID,
			&sm.CreatedByUserName,
			&sm.CancelledByUserID,
			&sm.CancelledByUserName,
		)
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, &sm)
	}

	return &result, nil
}

type CreateStockMovementParams struct {
	Type      string
	Date      time.Time
	EntityID  *pgxuuid.UUID
	CreatedBy *pgxuuid.UUID
}

func (r *PgRepository) CreateStockMovement(ctx context.Context, param *CreateStockMovementParams) (*StockMovement, error) {
	var sm StockMovement
	err := r.db.QueryRow(ctx, `
		INSERT INTO "stock_movements" (
			status,
			type,
			date,
			entity_id,
			created_by_user_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) RETURNING
			id,
			status,
			type,
			date,
			created_at,
			updated_at
	`, "ACTIVE", param.Type, param.Date, param.EntityID, param.CreatedBy).Scan(
		&sm.ID,
		&sm.Status,
		&sm.Type,
		&sm.Date,
		&sm.CreatedAt,
		&sm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sm, nil
}

type UpdateStockMovementParams struct {
	ID       *pgxuuid.UUID
	Date     time.Time
	EntityID *pgxuuid.UUID
}

func (r *PgRepository) UpdateStockMovement(ctx context.Context, param *UpdateStockMovementParams) (*StockMovement, error) {
	var sm StockMovement
	err := r.db.QueryRow(ctx, `
		UPDATE "stock_movements" SET
			date = $2,
			entity_id = $3
		WHERE
			id = $1
		RETURNING
			id,
			status,
			type,
			date,
			created_at,
			updated_at
	`, param.ID, param.Date, param.EntityID).Scan(
		&sm.ID,
		&sm.Status,
		&sm.Type,
		&sm.Date,
		&sm.CreatedAt,
		&sm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sm, nil
}

type DeleteStockMovementParams struct {
	ID     *pgxuuid.UUID
	UserID *pgxuuid.UUID
}

// Doest not delete the record, just changes the status to INACTIVE
func (r *PgRepository) DeleteStockMovement(ctx context.Context, param *DeleteStockMovementParams) (*StockMovement, error) {
	var sm StockMovement
	err := r.db.QueryRow(ctx, `
		UPDATE "stock_movements" SET
			status = 'INACTIVE',
			cancelled_by_user_id = $2,
			updated_at = now()
		WHERE
			id = $1
		RETURNING
			id,
			status,
			type,
			date,
			created_at,
			updated_at
	`, param.ID, param.UserID).Scan(
		&sm.ID,
		&sm.Status,
		&sm.Type,
		&sm.Date,
		&sm.CreatedAt,
		&sm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sm, nil
}

func (r *PgRepository) FetchStockMovementByID(ctx context.Context, id *pgxuuid.UUID) (*StockMovement, error) {
	var sm StockMovement
	err := r.db.QueryRow(ctx, `
		SELECT
			sm.id,
			sm.status,
			sm.type,
			sm.date,
			sm.created_at,
			sm.updated_at,
			e.id,
			coalesce(e.name, '-'),
			coalesce(e.ruc, e.ci, '-'),
			cu.id,
			cu.name,
			cu2.id,
			coalesce(cu2.name, '-')
		FROM "stock_movements" sm
		LEFT JOIN "entities" e ON e.id = sm.entity_id
		LEFT JOIN "users" cu ON cu.id = sm.created_by_user_id
		LEFT JOIN "users" cu2 ON cu2.id = sm.cancelled_by_user_id
		WHERE
			sm.id = $1
	`, id).Scan(
		&sm.ID,
		&sm.Status,
		&sm.Type,
		&sm.Date,
		&sm.CreatedAt,
		&sm.UpdatedAt,
		&sm.EntityID,
		&sm.EntityName,
		&sm.EntityDocument,
		&sm.CreatedByUserID,
		&sm.CreatedByUserName,
		&sm.CancelledByUserID,
		&sm.CancelledByUserName,
	)
	if err != nil {
		return nil, err
	}

	return &sm, nil
}
