package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/services"
)

func (h *Handlers) RegisterEntityRoutes() {
	g := h.app.Group("/entities")

	g.Get("/", h.getAllEntities)
	g.Get("/:id", h.getEntityById)
	g.Post("/", h.createEntity)
	g.Put("/:id", h.updateEntity)
}

type GetAllEntitiesQuery struct {
	StatusOptions []string `query:"status"`
	Search        string   `query:"search"`
	Limit         int      `query:"limit"`
	Offset        int      `query:"offset"`
}

func (h *Handlers) getAllEntities(c *fiber.Ctx) error {
	params := new(GetAllEntitiesQuery)
	if err := c.QueryParser(params); err != nil {
		return constants.InvalidBody()
	}

	if params.Limit <= 9 {
		params.Limit = 10
	}

	response, err := h.sm.FetchEntities(c.Context(), &services.FetchEntitiesParams{
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

func (h *Handlers) getEntityById(c *fiber.Ctx) error {
	entityId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	entity, err := h.sm.GetEntityByID(c.Context(), entityId)
	if err != nil {
		return err
	}

	if entity == nil {
		return c.Status(fiber.StatusNotFound).Send([]byte{})
	}

	return c.Status(fiber.StatusOK).JSON(entity)
}

type CreateEntityBody struct {
	Name string `json:"name"`
	RUC  string `json:"ruc"`
	CI   string `json:"ci"`
}

func (h *Handlers) createEntity(c *fiber.Ctx) error {
	body := new(CreateEntityBody)
	if err := c.BodyParser(body); err != nil {
		return constants.InvalidBody()
	}

	entity, err := h.sm.CreateEntity(c.Context(), &services.CreateEntityParams{
		Name: body.Name,
		RUC:  body.RUC,
		CI:   body.CI,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(entity)
}

type UpdateEntityBody struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	RUC    string `json:"ruc"`
	CI     string `json:"ci"`
}

func (h *Handlers) updateEntity(c *fiber.Ctx) error {
	entityId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	body := new(UpdateEntityBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	entity, err := h.sm.UpdateEntity(c.Context(), &services.UpdateEntityParams{
		ID:     entityId,
		Status: body.Status,
		Name:   body.Name,
		RUC:    body.RUC,
		CI:     body.CI,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(entity)
}
