package repository

import (
	"context"
	"fmt"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"math"
	"time"
)

type StockMovement struct {
	ID        *pgxuuid.UUID
	Status    string
	Type      string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Total int

	EntityID       *pgxuuid.UUID
	EntityName     string
	EntityDocument string

	CreatedByUserID     *pgxuuid.UUID
	CreatedByUserName   string
	CancelledByUserID   *pgxuuid.UUID
	CancelledByUserName string

	Items []*StockMovementItem
}

type StockMovementItem struct {
	ID              *pgxuuid.UUID
	StockMovementID *pgxuuid.UUID
	ProductID       *pgxuuid.UUID
	ProductName     string
	Quantity        int
	Price           int
	Total           int
	Batch           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
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
			coalesce(e.name, ''),
			coalesce(e.ruc, e.ci, ''),
			cu.id,
			cu.name,
			cu2.id,
			coalesce(cu2.name, '')
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
	defer rows.Close()

	var result FetchStockMovementsResult
	for rows.Next() {
		sm := StockMovement{}
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

	smIDs := make([]pgxuuid.UUID, len(result.Items))
	for i, sm := range result.Items {
		smIDs[i] = *sm.ID
	}

	rows, err = r.db.Query(ctx, `
		SELECT
			smi.id,
			smi.stock_movement_id,
			smi.product_id,
			p.name,
			smi.quantity,
			smi.price,
			smi.batch,
			smi.updated_at,
			smi.updated_at
		FROM "stock_movement_items" smi
		LEFT JOIN products p on p.id = smi.product_id
		WHERE
			smi.stock_movement_id = ANY($1::uuid[])`, smIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	smMap := make(map[pgxuuid.UUID]*StockMovement, len(result.Items))
	for _, sm := range result.Items {
		smMap[*sm.ID] = sm
	}

	for rows.Next() {
		smi := StockMovementItem{}
		err := rows.Scan(
			&smi.ID,
			&smi.StockMovementID,
			&smi.ProductID,
			&smi.ProductName,
			&smi.Quantity,
			&smi.Price,
			&smi.Batch,
			&smi.CreatedAt,
			&smi.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if sm, ok := smMap[*smi.StockMovementID]; ok {
			sm.Items = append(sm.Items, &smi)
			sm.Total += int(math.Round(float64(smi.Quantity*smi.Price) / 1000))
		}
	}

	return &result, nil
}

type CreateStockMovementParams struct {
	Type      string
	Date      time.Time
	EntityID  *pgxuuid.UUID
	CreatedBy *pgxuuid.UUID
	Items     []*CreateStockItem
}

type CreateStockItem struct {
	ProductID *pgxuuid.UUID
	Quantity  int
	Price     int
	Batch     string
}

func (r *PgRepository) CreateStockMovement(ctx context.Context, params *CreateStockMovementParams) (*StockMovement, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	var smID pgxuuid.UUID
	err = tx.QueryRow(ctx, `
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
		) RETURNING id
	`, "ACTIVE", params.Type, params.Date, params.EntityID, params.CreatedBy).Scan(
		&smID,
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("smID: %v\n", smID)

	for _, item := range params.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO "stock_movement_items" (
				stock_movement_id,
				product_id,
				quantity,
				price,
				batch
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			) RETURNING id
		`, smID, item.ProductID, item.Quantity, item.Price, item.Batch)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	stockMovement, err := r.FetchStockMovementByID(ctx, &smID)
	if err != nil {
		return nil, err
	}

	return stockMovement, nil
}

type UpdateStockMovementParams struct {
	ID       *pgxuuid.UUID
	Date     time.Time
	EntityID *pgxuuid.UUID
}

func (r *PgRepository) UpdateStockMovement(ctx context.Context, param *UpdateStockMovementParams) (*StockMovement, error) {
	var smID pgxuuid.UUID
	err := r.db.QueryRow(ctx, `
		UPDATE "stock_movements" SET
			date = $2,
			entity_id = $3
		WHERE
			id = $1
	`, param.ID, param.Date, param.EntityID).Scan(
		&smID,
	)
	if err != nil {
		return nil, err
	}

	stockMovement, err := r.FetchStockMovementByID(ctx, &smID)
	if err != nil {
		return nil, err
	}

	return stockMovement, nil
}

type DeleteStockMovementParams struct {
	ID     *pgxuuid.UUID
	UserID *pgxuuid.UUID
}

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
			coalesce(e.name, ''),
			coalesce(e.ruc, e.ci, ''),
			cu.id,
			cu.name,
			cu2.id,
			coalesce(cu2.name, '')
		FROM "stock_movements" sm
		LEFT JOIN "entities" e ON e.id = sm.entity_id
		LEFT JOIN "users" cu ON cu.id = sm.created_by_user_id
		LEFT JOIN "users" cu2 ON cu2.id = sm.cancelled_by_user_id
		WHERE
			sm.id = $1
	`, &id).Scan(
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

	rows, err := r.db.Query(ctx, `
		SELECT 
		    	smi.id,
		    	smi.product_id,
		    	p.name,
		    	smi.stock_movement_id,
		    	smi.quantity,
		    	smi.price,
		    	coalesce(batch, ''),
		    	smi.created_at,
		    	smi.updated_at
		FROM "stock_movement_items" smi
		LEFT JOIN "products" p on smi.product_id = p.id
		WHERE smi.stock_movement_id = $1
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total int

	sm.Items = make([]*StockMovementItem, 0)
	for rows.Next() {
		item := StockMovementItem{}
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.StockMovementID,
			&item.Quantity,
			&item.Price,
			&item.Batch,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("Error while scanning stock_movement_item")
			return nil, err
		}

		total = total + int(math.Round(float64(item.Quantity*item.Price)/1000))

		sm.Items = append(sm.Items, &item)
	}

	sm.Total = total

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &sm, nil
}
