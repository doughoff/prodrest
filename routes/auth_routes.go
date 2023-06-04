package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/middleware"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
)

func (h *Handlers) RegisterAuthRoutes() {
	// register routes here
	g := h.app.Group("/auth")

	g.Post("/login", h.login)
	g.Post("/logout", h.logout)
	g.Get("/me", h.profile)
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *Handlers) login(c *fiber.Ctx) error {
	payload := new(LoginPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	ok, user, err := h.sm.CheckEmailAndPassword(c.Context(), payload.Email, payload.Password)
	if err != nil {
		return err
	}

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "invalid email or password"})
	}

	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "invalid email or password"})
	}

	// generate new uuid for session key
	sessionId, err := uuid.NewV4()
	if err != nil {
		return err
	}

	sessionData := &middleware.SessionData{
		UserId: user.ID.String(),
		Roles:  user.Roles,
	}
	sessionDataBytes, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}
	err = h.sessionStore.Set(sessionId.String(), sessionDataBytes, constants.SessionDuration)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{"session_id": sessionId.String()})
}

func (h *Handlers) logout(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	sessionID := headers["X-Session"]
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "missing session id"})
	}

	err := h.sessionStore.Delete(sessionID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send([]byte{})
}

func (h *Handlers) profile(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	if userId == nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "missing user id"})
	}
	strid, ok := userId.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "invalid user id"})
	}
	userUUID, err := uuid.FromString(strid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "invalid user id"})
	}

	pgxUserUUID := pgxuuid.UUID(userUUID.Bytes())
	user, err := h.sm.FetchUserById(c.Context(), &pgxUserUUID)
	if err != nil {
		return err
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).Send([]byte{})
	}

	return c.Status(fiber.StatusOK).JSON(user)

}
