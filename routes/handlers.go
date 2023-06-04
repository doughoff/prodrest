package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/gofrs/uuid"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/services"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
)

type Handlers struct {
	sm           *services.ServiceManager
	app          *fiber.App
	sessionStore *memory.Storage
}

func NewHandlers(app *fiber.App, serviceManager *services.ServiceManager, store *memory.Storage) *Handlers {
	return &Handlers{
		sm:           serviceManager,
		app:          app,
		sessionStore: store,
	}
}

func (h *Handlers) getIdParam(c *fiber.Ctx) (*pgxuuid.UUID, error) {
	param := struct {
		ID uuid.UUID `params:"id"`
	}{}

	err := c.ParamsParser(&param) // "{"id": 111}"
	if err != nil {
		return nil, constants.InvalidParams("invalid id on url query string")
	}

	id := pgxuuid.UUID(param.ID.Bytes())
	return &id, nil

}
