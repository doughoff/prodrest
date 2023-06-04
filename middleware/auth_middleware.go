package middleware

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory"
	"github.com/hoffax/prodrest/constants"
)

type SessionData struct {
	UserId string
	Roles  []string
}

func AuthMiddleware(store *memory.Storage) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/auth/login" {
			return c.Next()
		}

		headers := c.GetReqHeaders()
		sessionID := headers["X-Session"]
		if sessionID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthenticated")
		}

		sessionDataBytes, err := store.Get(sessionID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if sessionDataBytes == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthenticated")
		}

		var sessionData SessionData

		if err = json.Unmarshal(sessionDataBytes, &sessionData); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to decode session data")
		}

		// Set locals
		c.Locals("userId", sessionData.UserId)
		c.Locals("roles", sessionData.Roles)

		if err = store.Set(sessionID, sessionDataBytes, constants.SessionDuration); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to refresh session")
		}

		return c.Next()
	}
}
