package services

import (
	"context"
	uuid "github.com/jackc/pgx-gofrs-uuid"
)

type RecipeDTO struct {
	RecipeID          *uuid.UUID            `json:"RecipeId"`
	GroupID           *uuid.UUID            `json:"groupId"`
	Name              string                `json:"name"`
	Description       string                `json:"description"`
	Status            string                `json:"status"`
	Revision          int                   `json:"revision"`
	IsCurrent         bool                  `json:"isCurrent"`
	CreatedByUserID   *uuid.UUID            `json:"createdByUserId"`
	CreatedByUserName string                `json:"createdByUserName"`
	CreatedAt         string                `json:"createdAt"`
	Ingredients       []RecipeIngredientDTO `json:"ingredients"`
}

type RecipeIngredientDTO struct {
	ID          *uuid.UUID `json:"id"`
	RecipeID    *uuid.UUID `json:"recipeId"`
	ProductID   *uuid.UUID `json:"productId"`
	ProductName string     `json:"productName"`
	Quantity    int        `json:"quantity"`
}

type CreateRecipeParams struct {
	Name        string                   `json:"name" validate:"required,min=3,max=255"`
	Ingredients []CreateRecipeIngredient `json:"ingredients" validate:"required,dive,required"`
}

type CreateRecipeIngredient struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
}

type UpdateRecipeParams struct {
	ID uuid.UUID `json:"id" validate:"required"`
	CreateRecipeParams
}

type FetchRecipeParams struct {
	Search        string   `json:"search"`
	StatusOptions []string `json:"statusOptions"`
	Limit         int      `json:"limit"`
	Offset        int      `json:"offset"`
}

type FetchRecipeResponse struct {
	Recipes []RecipeDTO `json:"recipes"`
	Total   int         `json:"total"`
}

func (s *ServiceManager) FetchRecipes(ctx context.Context, params *FetchRecipeParams) (*FetchRecipeResponse, error) {
	return &FetchRecipeResponse{
		Recipes: []RecipeDTO{},
		Total:   0,
	}, nil
}

func (s *ServiceManager) CreateRecipe(ctx context.Context, params *CreateRecipeParams) (*RecipeDTO, error) {
	return &RecipeDTO{}, nil
}

func (s *ServiceManager) UpdateRecipe(ctx context.Context, params *UpdateRecipeParams) (*RecipeDTO, error) {
	return &RecipeDTO{}, nil
}

func (s *ServiceManager) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *ServiceManager) FetchRecipe(ctx context.Context, id uuid.UUID) (*RecipeDTO, error) {
	return &RecipeDTO{}, nil
}
