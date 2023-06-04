package services

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type ProductDTO struct {
	ID               *uuid.UUID `json:"id"`
	Status           string     `json:"status"`
	Name             string     `json:"name"`
	Barcode          string     `json:"barcode"`
	Unit             string     `json:"unit"`
	BatchControl     bool       `json:"batchControl"`
	ConversionFactor int        `json:"conversionFactor"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

func (s *ServiceManager) toProductDTO(product *repository.Product) *ProductDTO {
	productId, err := s.parseUUID(product.ID)
	if err != nil {
		productId = nil
	}

	return &ProductDTO{
		ID:               productId,
		Status:           product.Status,
		Name:             product.Name,
		Barcode:          product.Barcode,
		Unit:             product.Unit,
		BatchControl:     product.BatchControl,
		ConversionFactor: product.ConversionFactor,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}
}

type CreateProductParams struct {
	Barcode          string `validate:"required,gte=3"`
	Name             string `validate:"required,gte=3,lte=80"`
	Unit             string `validate:"required,custom_unit"`
	ConversionFactor int    `validate:"required,gt=0"`
	BatchControl     bool
}

func (s *ServiceManager) CreateProduct(ctx context.Context, params *CreateProductParams) (*ProductDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.GetProductByBarcode(ctx, params.Barcode)
	if err == nil {
		return nil, NewUniqueConstrainError("barcode")
	} else {
		if err != pgx.ErrNoRows {
			return nil, err
		}
	}

	product, err := s.repo.CreateProduct(ctx, &repository.CreateProductParams{
		Barcode:          params.Barcode,
		Name:             params.Name,
		Unit:             params.Unit,
		BatchControl:     params.BatchControl,
		ConversionFactor: params.ConversionFactor,
	})
	if err != nil {
		return nil, err
	}

	return s.toProductDTO(product), nil
}

type UpdateProductParams struct {
	ID               *pgxuuid.UUID `validate:"required"`
	Status           string        `validate:"required,custom_status"`
	Barcode          string        `validate:"required,gte=3"`
	Name             string        `validate:"required,gte=3,lte=80"`
	Unit             string        `validate:"required,custom_unit"`
	BatchControl     bool
	ConversionFactor int `validate:"required,gte=1"`
}

func (s *ServiceManager) UpdateProduct(ctx context.Context, params *UpdateProductParams) (*ProductDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.GetProductByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	if product.Barcode != params.Barcode {
		_, err = s.repo.GetProductByBarcode(ctx, params.Barcode)
		if err == nil {
			return nil, NewUniqueConstrainError("barcode")
		} else {
			if err != pgx.ErrNoRows {
				return nil, err
			}
		}
	}

	product, err = s.repo.UpdateProduct(ctx, &repository.UpdateProductParams{
		ID:               params.ID,
		Status:           params.Status,
		Barcode:          params.Barcode,
		Name:             params.Name,
		Unit:             params.Unit,
		BatchControl:     params.BatchControl,
		ConversionFactor: params.ConversionFactor,
	})
	if err != nil {
		return nil, err
	}

	return s.toProductDTO(product), nil
}

type FetchProductsParams struct {
	Search        string
	StatusOptions []string `validate:"dive,custom_status"`
	Limit         int      `validate:"required,gte=1,lte=100"`
	Offset        int      `validate:"gte=0"`
}

type FetchProductsDTOResult struct {
	TotalCount int           `json:"totalCount"`
	Items      []*ProductDTO `json:"items"`
}

func (s *ServiceManager) FetchProducts(ctx context.Context, params *FetchProductsParams) (*FetchProductsDTOResult, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	result, err := s.repo.FetchProducts(ctx, &repository.FetchProductsParams{
		Search:        params.Search,
		StatusOptions: params.StatusOptions,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return nil, err
	}

	itemsDTP := make([]*ProductDTO, 0)
	for _, item := range result.Items {
		itemsDTP = append(itemsDTP, s.toProductDTO(item))
	}

	return &FetchProductsDTOResult{
		TotalCount: result.TotalCount,
		Items:      itemsDTP,
	}, nil
}

func (s *ServiceManager) FetchProductById(ctx context.Context, id *pgxuuid.UUID) (*ProductDTO, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return s.toProductDTO(product), nil
}

func (s *ServiceManager) FetchProductByBarcode(ctx context.Context, barcode string) (*ProductDTO, error) {
	product, err := s.repo.GetProductByBarcode(ctx, barcode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return s.toProductDTO(product), nil
}
