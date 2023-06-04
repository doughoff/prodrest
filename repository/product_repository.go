package repository

import (
	"context"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"time"
)

type Product struct {
	ID               *pgxuuid.UUID
	Status           string
	Name             string
	Barcode          string
	Unit             string
	BatchControl     bool
	ConversionFactor int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type FetchProductsParams struct {
	StatusOptions []string
	Search        string
	Limit         int
	Offset        int
}

type FetchProductsResult struct {
	TotalCount int
	Items      []*Product
}

func (r *PgRepository) FetchProducts(ctx context.Context, params *FetchProductsParams) (*FetchProductsResult, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
		    COUNT(*) OVER() AS full_count,
			id,
			status,
			name,
			barcode,
			unit,
			batch_control,
			conversion_factor,
			created_at,
			updated_at
		FROM "products"
		WHERE
		    status = ANY($1::status[])
			 OR (
			     name ILIKE '%' || $2 || '%'
			     OR barcode ILIKE '%' || $2 || '%'
			 )
		ORDER BY
		    created_at DESC
		LIMIT $3
		OFFSET $4
	`, params.StatusOptions, params.Search, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalCount int
	products := make([]*Product, 0)
	for rows.Next() {
		product := Product{}
		err := rows.Scan(
			&totalCount,
			&product.ID,
			&product.Status,
			&product.Name,
			&product.Barcode,
			&product.Unit,
			&product.BatchControl,
			&product.ConversionFactor,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return &FetchProductsResult{
		TotalCount: totalCount,
		Items:      products,
	}, nil
}

func (r *PgRepository) GetProductByID(ctx context.Context, id *pgxuuid.UUID) (*Product, error) {
	product := Product{}
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			status,
			name,
			barcode,
			unit,
			batch_control,
			conversion_factor,
			created_at,
			updated_at
		FROM "products"
		WHERE
			id = $1
	`, id).Scan(
		&product.ID,
		&product.Status,
		&product.Name,
		&product.Barcode,
		&product.Unit,
		&product.BatchControl,
		&product.ConversionFactor,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *PgRepository) GetProductByBarcode(ctx context.Context, barcode string) (*Product, error) {
	product := Product{}
	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			status,
			name,
			barcode,
			unit,
			batch_control,
			conversion_factor,
			created_at,
			updated_at
		FROM "products"
		WHERE
			barcode = $1
	`, barcode).Scan(
		&product.ID,
		&product.Status,
		&product.Name,
		&product.Barcode,
		&product.Unit,
		&product.BatchControl,
		&product.ConversionFactor,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

type CreateProductParams struct {
	Name             string
	Barcode          string
	Unit             string
	BatchControl     bool
	ConversionFactor int
}

func (r *PgRepository) CreateProduct(ctx context.Context, params *CreateProductParams) (*Product, error) {
	product := Product{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO "products" (
			name,
			barcode,
			unit,
			batch_control,
			conversion_factor
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING
			id,
			status,
			name,
			barcode,
			unit,
			batch_control,
			conversion_factor,
			created_at,
			updated_at
	`,
		params.Name,
		params.Barcode,
		params.Unit,
		params.BatchControl,
		params.ConversionFactor,
	).Scan(
		&product.ID,
		&product.Status,
		&product.Name,
		&product.Barcode,
		&product.Unit,
		&product.BatchControl,
		&product.ConversionFactor,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

type UpdateProductParams struct {
	ID               *pgxuuid.UUID
	Status           string
	Name             string
	Barcode          string
	Unit             string
	BatchControl     bool
	ConversionFactor int
}

func (r *PgRepository) UpdateProduct(ctx context.Context, params *UpdateProductParams) (*Product, error) {
	product := Product{}
	err := r.db.QueryRow(ctx, `
		UPDATE "products" SET
			status = $2,
			name = $3,
			barcode = $4,
			unit = $5,
			batch_control = $6,
			conversion_factor = $7
		WHERE
			id = $1
		RETURNING
			id,
			status,
			name,
			barcode,
			unit,
			batch_control,
			conversion_factor,
			created_at,
			updated_at
	`,
		params.ID,
		params.Status,
		params.Name,
		params.Barcode,
		params.Unit,
		params.BatchControl,
		params.ConversionFactor,
	).Scan(
		&product.ID,
		&product.Status,
		&product.Name,
		&product.Barcode,
		&product.Unit,
		&product.BatchControl,
		&product.ConversionFactor,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
