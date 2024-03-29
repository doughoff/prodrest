package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/services"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"time"
)

func (h *Handlers) RegisterStockMovementRoutes() {
	g := h.app.Group("/stock_movements")

	g.Get("/", h.getAllStockMovements)
	g.Get("/:id", h.getStockMovementById)
	g.Post("/", h.createStockMovement)
	g.Put("/:id", h.updateStockMovement)
	g.Delete("/:id", h.CancelStockMovementByID)
}

type GetAllStockMovementsQuery struct {
	StatusOptions []string `query:"status"`
	TypeOptions   []string `query:"type"`
	StartDate     string   `query:"startDate"`
	Limit         int      `query:"limit"`
	Offset        int      `query:"offset"`
}

func (h *Handlers) getAllStockMovements(c *fiber.Ctx) error {
	params := new(GetAllStockMovementsQuery)
	if err := c.QueryParser(params); err != nil {
		return err
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, params.StartDate)
	if err != nil {
		return constants.InvalidParams("invalid startDate format")
	}

	stockMovements, err := h.sm.FetchStockMovements(c.Context(), &services.FetchStockMovementsParams{
		StatusOptions: params.StatusOptions,
		TypeOptions:   params.TypeOptions,
		StartDate:     startDate,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(stockMovements)
}

func (h *Handlers) getStockMovementById(c *fiber.Ctx) error {
	stockMovementId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	stockMovement, err := h.sm.FetchStockMovementByID(c.Context(), stockMovementId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(stockMovement)
}

type CreateStockMovementBody struct {
	Type     string         `json:"type"`
	Date     string         `json:"date"`
	EntityId uuid.UUID      `json:"entityId"`
	Items    []*CreateItems `json:"items"`
}

type CreateItems struct {
	ProductID *uuid.UUID `json:"productId"`
	Quantity  int        `json:"quantity"`
	Price     int        `json:"price"`
	Batch     string     `json:"batch"`
}

func (h *Handlers) createStockMovement(c *fiber.Ctx) error {
	params := new(CreateStockMovementBody)
	if err := c.BodyParser(params); err != nil {
		return constants.InvalidBody()
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, params.Date)
	if err != nil {
		return constants.NewRequiredFieldError("date")
	}

	userIDStr, ok := c.Locals("userId").(string)
	if !ok {
		return fiber.ErrUnauthorized
	}

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	items := make([]*services.CreateStockItem, 0)
	for _, item := range params.Items {
		itemUUID := pgxuuid.UUID(item.ProductID.Bytes())
		items = append(items, &services.CreateStockItem{
			ProductID: &itemUUID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Batch:     item.Batch,
		})
	}

	pgxUserID := pgxuuid.UUID(userID.Bytes())
	entityID := pgxuuid.UUID(params.EntityId.Bytes())
	stockMovement, err := h.sm.CreateStockMovement(c.Context(), &services.CreateStockMovementParams{
		Type:     params.Type,
		Date:     date,
		EntityID: &entityID,
		UserID:   &pgxUserID,
		Items:    items,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(stockMovement)
}

type UpdateStockMovementBody struct {
	Date     string    `validate:"required"`
	EntityID uuid.UUID `validate:"required"`
}

func (h *Handlers) updateStockMovement(c *fiber.Ctx) error {
	stockMovementId, err := h.getIdParam(c)
	if err != nil {
		return nil
	}

	params := new(UpdateStockMovementBody)
	if err := c.BodyParser(params); err != nil {
		return constants.InvalidBody()
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, params.Date)
	if err != nil {
		return constants.NewRequiredFieldError("date")
	}

	entityID := pgxuuid.UUID(params.EntityID.Bytes())
	stockMovement, err := h.sm.UpdateStockMovement(c.Context(), &services.UpdateStockMovementParams{
		ID:       stockMovementId,
		Date:     date,
		EntityID: &entityID,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(stockMovement)
}

func (h *Handlers) CancelStockMovementByID(c *fiber.Ctx) error {
	stockMovementId, err := h.getIdParam(c)
	if err != nil {
		return nil
	}

	userIDStr, ok := c.Locals("userId").(string)
	if !ok {
		fmt.Printf("userIdStr: %v\n", userIDStr)
		return fiber.ErrUnauthorized
	}

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fiber.ErrUnauthorized
	}

	pgxUserID := pgxuuid.UUID(userID.Bytes())
	stockMovement, err := h.sm.CancelStockMovementByID(c.Context(), &services.CancelStockMovementParams{
		ID:     stockMovementId,
		UserID: &pgxUserID,
	})

	return c.Status(fiber.StatusOK).JSON(stockMovement)
}
