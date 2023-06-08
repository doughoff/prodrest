package services

import (
	"bytes"
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"time"
)

type StockMovementDTO struct {
	ID        *uuid.UUID `json:"id"`
	Status    string     `json:"status"`
	Type      string     `json:"type"`
	Date      time.Time  `json:"date"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`

	EntityID            *uuid.UUID `json:"entityId"`
	EntityName          string     `json:"entityName"`
	EntityDocument      string     `json:"entityDocument"`
	CreatedByUserID     *uuid.UUID `json:"createdByUserId"`
	CreatedByUserName   string     `json:"createdByUserName"`
	CancelledByUserID   *uuid.UUID `json:"cancelledByUserId"`
	CancelledByUserName string     `json:"cancelledByUserName"`
}

func (s *ServiceManager) toStockMovementDTO(stockMovement *repository.StockMovement) *StockMovementDTO {
	stockMovementId, err := s.parseUUID(stockMovement.ID)
	if err != nil {
		stockMovementId = nil
	}

	entityId, err := s.parseUUID(stockMovement.EntityID)
	if err != nil {
		entityId = nil
	}

	createdByUserId, err := s.parseUUID(stockMovement.CreatedByUserID)
	if err != nil {
		createdByUserId = nil
	}

	cancelledByUserId, err := s.parseUUID(stockMovement.CancelledByUserID)
	if err != nil {
		cancelledByUserId = nil
	}

	return &StockMovementDTO{
		ID:                  stockMovementId,
		Status:              stockMovement.Status,
		Type:                stockMovement.Type,
		Date:                stockMovement.Date,
		CreatedAt:           stockMovement.CreatedAt,
		UpdatedAt:           stockMovement.UpdatedAt,
		EntityID:            entityId,
		EntityName:          stockMovement.EntityName,
		EntityDocument:      stockMovement.EntityDocument,
		CreatedByUserID:     createdByUserId,
		CreatedByUserName:   stockMovement.CreatedByUserName,
		CancelledByUserID:   cancelledByUserId,
		CancelledByUserName: stockMovement.CancelledByUserName,
	}
}

type CreateStockMovementParams struct {
	Type     string    `validate:"required"`
	Date     time.Time `validate:"required"`
	EntityID *pgxuuid.UUID
	UserID   *pgxuuid.UUID `validate:"required"`
}

func (s *ServiceManager) CreateStockMovement(ctx context.Context, params *CreateStockMovementParams) (*StockMovementDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	if params.Type == "PURCHASE" || params.Type == "SALE" {
		entityUUID, err := params.EntityID.UUIDValue()
		if err != nil {
			return nil, constants.NewRequiredFieldError("entityId")
		}
		if bytes.Equal(entityUUID.Bytes[:], uuid.Nil.Bytes()) {
			return nil, constants.NewRequiredFieldError("entityId")
		}
	} else {
		params.EntityID = nil
	}

	stockMovement, err := s.repo.CreateStockMovement(ctx, &repository.CreateStockMovementParams{
		Type:      params.Type,
		Date:      params.Date,
		EntityID:  params.EntityID,
		CreatedBy: params.UserID,
	})
	if err != nil {
		return nil, err
	}

	return s.toStockMovementDTO(stockMovement), nil
}

type UpdateStockMovementParams struct {
	ID       *pgxuuid.UUID `validate:"required"`
	Date     time.Time     `validate:"required"`
	EntityID *pgxuuid.UUID
}

func (s *ServiceManager) UpdateStockMovement(ctx context.Context, params *UpdateStockMovementParams) (*StockMovementDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	currentStockMovement, err := s.repo.FetchStockMovementByID(ctx, params.ID)
	if err != nil {
		return nil, constants.NewNotFoundError()
	}

	if currentStockMovement.Status == "INACTIVE" {
		return nil, constants.NewInvalidOperationError("stock movement is inactive")
	}

	if currentStockMovement.Type == "PURCHASE" || currentStockMovement.Type == "SALE" {
		entityUUID, err := params.EntityID.UUIDValue()
		if err != nil {
			return nil, constants.NewRequiredFieldError("entityId")
		}
		if bytes.Equal(entityUUID.Bytes[:], uuid.Nil.Bytes()) {
			return nil, constants.NewRequiredFieldError("entityId")
		}
	} else {
		params.EntityID = nil
	}

	stockMovement, err := s.repo.UpdateStockMovement(ctx, &repository.UpdateStockMovementParams{
		ID:       params.ID,
		Date:     params.Date,
		EntityID: params.EntityID,
	})
	if err != nil {
		return nil, err
	}

	return s.toStockMovementDTO(stockMovement), nil
}

type CancelStockMovementParams struct {
	ID     *pgxuuid.UUID `validate:"required"`
	UserID *pgxuuid.UUID `validate:"required"`
}

func (s *ServiceManager) CancelStockMovementByID(ctx context.Context, params *CancelStockMovementParams) (*StockMovementDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	currentStockMovement, err := s.repo.FetchStockMovementByID(ctx, params.ID)
	if err != nil {
		return nil, constants.NewNotFoundError()
	}

	if currentStockMovement.Status == "INACTIVE" {
		return nil, constants.NewInvalidOperationError("stock movement is inactive")
	}

	stockMovement, err := s.repo.DeleteStockMovement(ctx, &repository.DeleteStockMovementParams{
		ID:     params.ID,
		UserID: params.UserID,
	})
	if err != nil {
		return nil, err
	}

	return s.toStockMovementDTO(stockMovement), nil
}

type FetchStockMovementByIDParams struct {
	ID *pgxuuid.UUID `validate:"required"`
}

func (s *ServiceManager) FetchStockMovementByID(ctx context.Context, id *pgxuuid.UUID) (*StockMovementDTO, error) {
	stockMovement, err := s.repo.FetchStockMovementByID(ctx, id)
	if err != nil {
		return nil, constants.NewNotFoundError()
	}

	return s.toStockMovementDTO(stockMovement), nil
}

type FetchStockMovementsParams struct {
	StatusOptions []string  `validate:"dive,custom_status"`
	TypeOptions   []string  `validate:"dive,required"`
	StartDate     time.Time `validate:"required"`
	Limit         int       `validate:"required,gte=10,max=100"`
	Offset        int       `validate:"gte=0"`
}

type FetchStockMovementsResult struct {
	TotalCount int                 `json:"totalCount"`
	Items      []*StockMovementDTO `json:"items"`
}

func (s *ServiceManager) FetchStockMovements(ctx context.Context, params *FetchStockMovementsParams) (*FetchStockMovementsResult, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	stockMovements, err := s.repo.FetchStockMovements(ctx, &repository.FetchStockMovementsParams{
		StatusOptions: params.StatusOptions,
		TypeOptions:   params.TypeOptions,
		StartDate:     params.StartDate,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return nil, err
	}

	var stockMovementsDTO []*StockMovementDTO
	for _, stockMovement := range stockMovements.Items {
		stockMovementsDTO = append(stockMovementsDTO, s.toStockMovementDTO(stockMovement))
	}

	return &FetchStockMovementsResult{
		TotalCount: stockMovements.TotalCount,
		Items:      stockMovementsDTO,
	}, nil
}
