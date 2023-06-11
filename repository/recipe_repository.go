package repository

import (
	"context"
	"fmt"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
)

type Recipe struct {
	RecipeID          *pgxuuid.UUID
	GroupID           *pgxuuid.UUID
	Name              string
	Description       string
	Status            string
	Revision          int
	IsCurrent         bool
	CreatedByUserID   *pgxuuid.UUID
	CreatedByUserName string
	CreatedAt         string
	Ingredients       []*RecipeIngredient
}

type RecipeIngredient struct {
	ID          *pgxuuid.UUID
	RecipeID    *pgxuuid.UUID
	ProductID   *pgxuuid.UUID
	ProductName string
	Quantity    int
}

type FetchRecipeParams struct {
	Search        string
	StatusOptions []string
	Limit         int
	Offset        int
}

type FetchRecipeResult struct {
	TotalCount int
	Items      []*Recipe
}

type CreateRecipeParams struct {
	Name            string
	CreatedByUserID *pgxuuid.UUID
	Ingredients     []*CreateRecipeIngredient
}

type CreateRecipeIngredient struct {
	ProductID *pgxuuid.UUID
	Quantity  int
}

type UpdateRecipeParams struct {
	ID *pgxuuid.UUID
	CreateRecipeParams
}

func (r *PgRepository) FetchRecipes(ctx context.Context, params *FetchRecipeParams) (*FetchRecipeResult, error) {
	rows, err := r.db.Query(ctx, `
	select 
	    count(*) over() as full_count,
	    r.recipe_id,
	    r.recipe_group_id,
	    r.name,
	    r.status,
	    r.revision,
	    r.is_current,
	    r.created_by_user_id,
	    u.name,
	    r.created_at
	from "recipes" r
	left join "users" u on r.created_by_user_id = u.id
	where r.status = any($1::status[])
	  and r.name ILIKE '%' || $2 || '%'
	  and r.is_current = true
	order by r.created_at
	limit $4
	offset $5
`, params.StatusOptions, params.Search, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	result := &FetchRecipeResult{
		Items: make([]*Recipe, 0),
	}
	for rows.Next() {
		recipe := &Recipe{}
		err := rows.Scan(
			&result.TotalCount,
			&recipe.RecipeID,
			&recipe.GroupID,
			&recipe.Name,
			&recipe.Status,
			&recipe.Revision,
			&recipe.IsCurrent,
			&recipe.CreatedByUserID,
			&recipe.CreatedByUserName,
			&recipe.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, recipe)
	}

	result.TotalCount = len(result.Items)

	return result, nil
}

func (r *PgRepository) CreateRecipe(ctx context.Context, params *CreateRecipeParams) (*Recipe, error) {
	var createdUUID *pgxuuid.UUID
	tx, err := r.db.Begin(ctx)
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("error rolling back transaction: %v", err)
		}
	}(tx, ctx)

	err = tx.QueryRow(ctx, `
		insert into "recipes" (name, created_by_user_id)
		values ( $1, $2)
		returning recipe_id
	`, &params.Name, &params.CreatedByUserID).Scan(&createdUUID)
	if err != nil {
		return nil, err
	}

	for _, ingredient := range params.Ingredients {
		_, err := tx.Exec(ctx, `
			insert into "recipe_ingredients" (recipe_id, product_id, quantity)
			values ( $1, $2, $3)
		`, createdUUID, ingredient.ProductID, ingredient.Quantity)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	newRecipe, err := r.GetRecipeByID(ctx, createdUUID)
	if err != nil {
		return nil, err
	}

	return newRecipe, nil
}

func (r *PgRepository) UpdateRecipe(ctx context.Context, params *UpdateRecipeParams) (*Recipe, error) {
	recipe, err := r.GetRecipeByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("error while making transaction rollback")
		}
	}(tx, ctx)
	_, err = tx.Exec(ctx, `
		update recipes
			set is_current = false
		where 
		    recipe_id = $1
	`, &params.ID)
	if err != nil {
		return nil, err
	}

	var createdUUID *pgxuuid.UUID
	err = tx.QueryRow(ctx, `
		insert into "recipes" (name,recipe_group_id, created_by_user_id,revision)
		values ( $1, $2, $3, $4)
		returning recipe_id
	`, &params.Name, &recipe.GroupID, &params.CreatedByUserID, recipe.Revision+1).Scan(&createdUUID)
	if err != nil {
		return nil, err
	}

	for _, ingredient := range params.Ingredients {
		_, err := tx.Exec(ctx, `
			insert into "recipe_ingredients" (recipe_id, product_id, quantity)
			values ( $1, $2, $3)
		`, createdUUID, ingredient.ProductID, ingredient.Quantity)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	newRecipe, err := r.GetRecipeByID(ctx, createdUUID)

	return newRecipe, nil
}

func (r *PgRepository) GetRecipeByID(ctx context.Context, id *pgxuuid.UUID) (*Recipe, error) {
	recipe := &Recipe{
		Ingredients: make([]*RecipeIngredient, 0),
	}
	err := r.db.QueryRow(ctx, `
		select 
			r.recipe_id,
			r.recipe_group_id,
			r.name,
			r.status,
			r.revision,
			r.is_current,
			r.created_by_user_id,
			u.name,
			r.created_at
		from "recipes" r
		left join "users" u on r.created_by_user_id = u.id
		where r.recipe_id = $1
	`, id).Scan(
		&recipe.RecipeID,
		&recipe.GroupID,
		&recipe.Name,
		&recipe.Status,
		&recipe.Revision,
		&recipe.IsCurrent,
		&recipe.CreatedByUserID,
		&recipe.CreatedByUserName,
		&recipe.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
	select 
	    ri.id,
	    ri.product_id,
	    p.name,
	    ri.quantity
from "recipe_ingredients" ri 
	left join "products" p on p.id = ri.product_id
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ingredient := &RecipeIngredient{}
		err = rows.Scan(
			&ingredient.ID,
			&ingredient.ProductID,
			&ingredient.ProductName,
			&ingredient.Quantity,
		)
		if err != nil {
			return nil, err
		}

		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	return recipe, nil
}
