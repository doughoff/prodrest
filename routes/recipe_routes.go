package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/services"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
)

func (h *Handlers) RegisterRecipeRoutes() {
	g := h.app.Group("/recipes")

	g.Get("/", h.getAllRecipes)
	g.Post("/", h.createRecipe)
	g.Put("/", h.updateRecipe)
	g.Delete("/:id", h.deleteRecipeById)
	g.Get("/:id", h.getRecipeById)
}

type GetAllRecipesQuery struct {
	StatusOptions []string `query:"status"`
	Search        string   `query:"search"`
	Limit         int      `query:"limit"`
	Offset        int      `query:"offset"`
}

func (h *Handlers) getAllRecipes(c *fiber.Ctx) error {
	params := new(GetAllRecipesQuery)
	if err := c.QueryParser(params); err != nil {
		return constants.InvalidBody()
	}

	if params.Limit <= 9 {
		params.Limit = 10
	}

	fmt.Printf("params: %v", params)

	response, err := h.sm.FetchRecipes(c.Context(), &services.FetchRecipeParams{
		StatusOptions: params.StatusOptions,
		Search:        params.Search,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

type CreateRecipe struct {
	Name        string                   `json:"name" validate:"required,min=3,max=255"`
	Ingredients []CreateRecipeIngredient `json:"ingredients" validate:"required,dive,required"`
}

type CreateRecipeIngredient struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
}

func (h *Handlers) createRecipe(c *fiber.Ctx) error {
	body := new(CreateRecipe)
	if err := c.BodyParser(body); err != nil {
		return constants.InvalidBody()
	}

	parsedIngredients := make([]*services.CreateRecipeIngredient, len(body.Ingredients))
	for i, ingredient := range body.Ingredients {
		productID := pgxuuid.UUID(ingredient.ProductID.Bytes())
		parsedIngredients[i] = &services.CreateRecipeIngredient{
			ProductID: &productID,
			Quantity:  ingredient.Quantity,
		}
	}

	response, err := h.sm.CreateRecipe(c.Context(), &services.CreateRecipeParams{
		Name: body.Name,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

type UpdateRecipe struct {
	CreateRecipe
}

func (h *Handlers) updateRecipe(c *fiber.Ctx) error {
	recipeId, err := h.getIdParam(c)
	if err != nil {
		return err
	}
	body := new(UpdateRecipe)
	if err := c.BodyParser(body); err != nil {
		return constants.InvalidBody()
	}

	recipeIngredients := make([]*services.CreateRecipeIngredient, 0)
	for _, item := range body.Ingredients {
		productID := pgxuuid.UUID(item.ProductID.Bytes())
		recipeIngredients = append(recipeIngredients, &services.CreateRecipeIngredient{
			ProductID: &productID,
			Quantity:  item.Quantity,
		})
	}

	createRecipeParams := &services.CreateRecipeParams{
		Name:        body.Name,
		Ingredients: recipeIngredients,
	}

	response, err := h.sm.UpdateRecipe(c.Context(), &services.UpdateRecipeParams{
		ID:                 recipeId,
		CreateRecipeParams: *createRecipeParams,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handlers) deleteRecipeById(c *fiber.Ctx) error {
	recipeId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	//response, err := h.sm.DeleteRecipe(c.Context(), &services.DeleteRecipeParams{
	//	ID: entityId,
	//})
	//if err != nil {
	//	return uuid.UUID{}, err
	//}

	return c.Status(fiber.StatusOK).JSON(recipeId)
}

func (h *Handlers) getRecipeById(c *fiber.Ctx) error {
	entityId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	recipe, err := h.sm.FetchRecipeByID(c.Context(), entityId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(recipe)
}
