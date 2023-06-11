package services

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
)

type RecipeDTO struct {
	RecipeID          *uuid.UUID             `json:"RecipeId"`
	GroupID           *uuid.UUID             `json:"groupId"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Status            string                 `json:"status"`
	Revision          int                    `json:"revision"`
	IsCurrent         bool                   `json:"isCurrent"`
	CreatedByUserID   *uuid.UUID             `json:"createdByUserId"`
	CreatedByUserName string                 `json:"createdByUserName"`
	CreatedAt         string                 `json:"createdAt"`
	Ingredients       []*RecipeIngredientDTO `json:"ingredients"`
}

type RecipeIngredientDTO struct {
	ID          *uuid.UUID `json:"id"`
	RecipeID    *uuid.UUID `json:"recipeId"`
	ProductID   *uuid.UUID `json:"productId"`
	ProductName string     `json:"productName"`
	Quantity    int        `json:"quantity"`
}

type CreateRecipeParams struct {
	Name        string                    `json:"name" validate:"required,min=3,max=255"`
	UserID      *pgxuuid.UUID             `validate:"required"`
	Ingredients []*CreateRecipeIngredient `json:"ingredients" validate:"required,dive,required"`
}

type CreateRecipeIngredient struct {
	ProductID *pgxuuid.UUID `json:"product_id" validate:"required"`
	Quantity  int           `json:"quantity" validate:"required"`
}

type UpdateRecipeParams struct {
	ID *pgxuuid.UUID `json:"id" validate:"required"`
	CreateRecipeParams
}

type FetchRecipeParams struct {
	Search        string   `json:"search"`
	StatusOptions []string `json:"statusOptions"`
	Limit         int      `json:"limit"`
	Offset        int      `json:"offset"`
}

type FetchRecipeResponse struct {
	Recipes    []*RecipeDTO `json:"recipes"`
	TotalCount int          `json:"total"`
}

func (s *ServiceManager) toRecipeDTO(recipe *repository.Recipe) *RecipeDTO {
	recipeID, err := s.parseUUID(recipe.RecipeID)
	if err != nil {
		recipeID = nil
	}

	groupID, err := s.parseUUID(recipe.GroupID)
	if err != nil {
		groupID = nil
	}

	createdByUserID, err := s.parseUUID(recipe.CreatedByUserID)
	if err != nil {
		createdByUserID = nil
	}

	recipeDTO := &RecipeDTO{
		RecipeID:          recipeID,
		GroupID:           groupID,
		Status:            recipe.Status,
		Name:              recipe.Name,
		Revision:          recipe.Revision,
		IsCurrent:         recipe.IsCurrent,
		CreatedByUserID:   createdByUserID,
		CreatedByUserName: recipe.CreatedByUserName,
		CreatedAt:         recipe.CreatedAt,
		Ingredients:       make([]*RecipeIngredientDTO, 0),
	}

	for _, ingredient := range recipe.Ingredients {
		ingredientID, err := s.parseUUID(ingredient.ID)
		if err != nil {
			ingredientID = nil
		}

		productID, err := s.parseUUID(ingredient.ProductID)
		if err != nil {
			productID = nil
		}

		recipeDTO.Ingredients = append(recipeDTO.Ingredients, &RecipeIngredientDTO{
			ID:          ingredientID,
			RecipeID:    recipeID,
			ProductID:   productID,
			ProductName: ingredient.ProductName,
			Quantity:    ingredient.Quantity,
		})
	}

	return recipeDTO
}

func (s *ServiceManager) FetchRecipes(ctx context.Context, params *FetchRecipeParams) (*FetchRecipeResponse, error) {
	result, err := s.repo.FetchRecipes(ctx, &repository.FetchRecipeParams{
		StatusOptions: params.StatusOptions,
		Search:        params.Search,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return nil, err
	}

	recipes := make([]*RecipeDTO, 0)
	for _, item := range result.Items {
		recipes = append(recipes, s.toRecipeDTO(item))
	}

	return &FetchRecipeResponse{
		Recipes:    recipes,
		TotalCount: result.TotalCount,
	}, nil
}

func (s *ServiceManager) CreateRecipe(ctx context.Context, params *CreateRecipeParams) (*RecipeDTO, error) {
	var ingredients []*repository.CreateRecipeIngredient
	for _, ingredient := range params.Ingredients {

		ingredients = append(ingredients, &repository.CreateRecipeIngredient{
			ProductID: ingredient.ProductID,
			Quantity:  ingredient.Quantity,
		})
	}

	newRecipe, err := s.repo.CreateRecipe(ctx, &repository.CreateRecipeParams{
		Name:            params.Name,
		CreatedByUserID: params.UserID,
		Ingredients:     ingredients,
	})
	if err != nil {
		return nil, err
	}

	return s.toRecipeDTO(newRecipe), nil
}

func (s *ServiceManager) UpdateRecipe(ctx context.Context, params *UpdateRecipeParams) (*RecipeDTO, error) {
	var ingredients []*repository.CreateRecipeIngredient
	for _, ingredient := range params.Ingredients {

		ingredients = append(ingredients, &repository.CreateRecipeIngredient{
			ProductID: ingredient.ProductID,
			Quantity:  ingredient.Quantity,
		})
	}

	createRecipeParams := &repository.CreateRecipeParams{
		Name:            params.Name,
		CreatedByUserID: params.UserID,
		Ingredients:     ingredients,
	}

	newRecipe, err := s.repo.UpdateRecipe(ctx, &repository.UpdateRecipeParams{
		ID:                 params.ID,
		CreateRecipeParams: *createRecipeParams,
	})
	if err != nil {
		return nil, err
	}

	return s.toRecipeDTO(newRecipe), nil
}

func (s *ServiceManager) DeleteRecipe(ctx context.Context, id uuid.UUID) error {

	return nil
}

func (s *ServiceManager) FetchRecipe(ctx context.Context, id uuid.UUID) (*RecipeDTO, error) {
	return &RecipeDTO{}, nil
}
