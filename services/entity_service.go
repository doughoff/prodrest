package services

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type EntityDTO struct {
	ID        *uuid.UUID `json:"id"`
	Status    string     `json:"status"`
	Name      string     `json:"name"`
	RUC       string     `json:"ruc"`
	CI        string     `json:"ci"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (s *ServiceManager) toEntityDTO(entity *repository.Entity) *EntityDTO {
	entityId, err := s.parseUUID(entity.ID)
	if err != nil {
		entityId = nil
	}

	return &EntityDTO{
		ID:        entityId,
		Status:    entity.Status,
		Name:      entity.Name,
		RUC:       entity.RUC,
		CI:        entity.CI,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// helper function to validate RUC or CI
// if both are empty, return false
func (s *ServiceManager) validateRUCOrCI(ruc string, ci string) bool {
	if ruc == "" && ci == "" {
		return false
	}

	return true
}

type CreateEntityParams struct {
	Name string `validate:"required,gte=3,lte=80"`
	RUC  string
	CI   string
}

func (s *ServiceManager) CreateEntity(ctx context.Context, params *CreateEntityParams) (*EntityDTO, error) {
	if !s.validateRUCOrCI(params.RUC, params.CI) {
		return nil, constants.NewRequiredFieldError("ruc or ci, at least one is required")
	}

	if params.RUC != "" {
		_, err := s.repo.GetEntityByRUC(ctx, params.RUC)
		if err == nil {
			return nil, constants.NewUniqueConstrainError("ruc")
		} else {
			if err != pgx.ErrNoRows {
				return nil, err
			}
		}
	}

	if params.CI != "" {
		_, err := s.repo.GetEntityByCI(ctx, params.CI)
		if err == nil {
			return nil, constants.NewUniqueConstrainError("ci")
		} else {
			if err != pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}
	}

	entity, err := s.repo.CreateEntity(ctx, &repository.CreateEntityParams{
		Name: params.Name,
		RUC:  params.RUC,
		CI:   params.CI,
	})
	if err != nil {
		return nil, err
	}

	return s.toEntityDTO(entity), nil
}

type UpdateEntityParams struct {
	ID     *pgxuuid.UUID `validate:"required"`
	Status string        `validate:"required,custom_status"`
	Name   string        `validate:"required,gte=3,lte=80"`
	RUC    string
	CI     string
}

func (s *ServiceManager) UpdateEntity(ctx context.Context, params *UpdateEntityParams) (*EntityDTO, error) {
	entity, err := s.repo.GetEntityById(ctx, params.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, constants.NewNotFoundError()
		}
		return nil, err
	}

	if !s.validateRUCOrCI(params.RUC, params.CI) {
		return nil, constants.NewRequiredFieldError("ruc or ci, at least one is required")
	}

	if params.RUC != "" {
		entity, err := s.repo.GetEntityByRUC(ctx, params.RUC)
		if err == nil {
			if entity.ID != params.ID {
				return nil, constants.NewUniqueConstrainError("ruc")
			}
		} else {
			if err != pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}
	}

	if params.CI != "" {
		entity, err := s.repo.GetEntityByCI(ctx, params.CI)
		if err == nil {
			if entity.ID != params.ID {
				return nil, constants.NewUniqueConstrainError("ci")
			}
		} else {
			if err != pgx.ErrNoRows {
				return nil, err
			}
			return nil, err
		}
	}

	entity, err = s.repo.UpdateEntity(ctx, &repository.UpdateEntityParams{
		ID:     params.ID,
		Status: params.Status,
		Name:   params.Name,
		RUC:    params.RUC,
		CI:     params.CI,
	})
	if err != nil {
		return nil, err
	}

	return s.toEntityDTO(entity), nil
}

func (s *ServiceManager) GetEntityByID(ctx context.Context, id *pgxuuid.UUID) (*EntityDTO, error) {
	entity, err := s.repo.GetEntityById(ctx, id)
	if err != nil {
		return nil, constants.NewNotFoundError()
	}

	return s.toEntityDTO(entity), nil
}

type FetchEntitiesParams struct {
	StatusOptions []string `validate:"dive,custom_status"`
	Search        string   `validate:"required"`
	Limit         int      `validate:"required,gte=1,lte=100"`
	Offset        int      `validate:"required,gte=0"`
}

type FetchEntitiesResponse struct {
	TotalCount int          `json:"totalCount"`
	Items      []*EntityDTO `json:"items"`
}

func (s *ServiceManager) FetchEntities(ctx context.Context, params *FetchEntitiesParams) (*FetchEntitiesResponse, error) {
	result, err := s.repo.FetchEntities(ctx, &repository.FetchEntitiesParams{
		StatusOptions: params.StatusOptions,
		Search:        params.Search,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return nil, err
	}

	entitiesDTO := make([]*EntityDTO, 0)
	for _, entity := range result.Items {
		entitiesDTO = append(entitiesDTO, s.toEntityDTO(entity))
	}

	return &FetchEntitiesResponse{
		TotalCount: result.TotalCount,
		Items:      entitiesDTO,
	}, nil
}
