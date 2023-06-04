package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodrest/services"
)

func (h *Handlers) RegisterUserRoutes() {
	g := h.app.Group("/users")

	g.Get("/", h.getAll)
	g.Get("/:id", h.getById)
	g.Post("/", h.createUser)
	g.Put("/:id", h.updateUser)

	// email checker
	h.app.Get("/check_email/:email", h.getByEmail)
}

type GetAllUsersQuery struct {
	StatusOptions []string `query:"status"`
}

func (h *Handlers) getAll(c *fiber.Ctx) error {
	params := new(GetAllUsersQuery)
	fmt.Printf("/users user:%v\n", c.Locals("user"))
	if err := c.QueryParser(params); err != nil {
		return err
	}

	users, err := h.sm.FetchUsers(c.Context(), &services.FetchUsersParams{
		StatusOptions: params.StatusOptions,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *Handlers) getById(c *fiber.Ctx) error {
	userId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	user, err := h.sm.FetchUserById(c.Context(), userId)
	if err != nil {
		return err
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).Send([]byte{})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *Handlers) getByEmail(c *fiber.Ctx) error {
	param := struct {
		Email string `params:"email" validate:"email,required"`
	}{}

	err := c.ParamsParser(&param)
	if err != nil {
		return h.InvalidParams("invalid email on url params")
	}

	user, err := h.sm.FetchUserByEmail(c.Context(), param.Email)
	if err != nil {
		return err
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).Send([]byte{})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

type CreateUserBody struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

func (h *Handlers) createUser(c *fiber.Ctx) error {
	userBody := new(CreateUserBody)

	if err := c.BodyParser(userBody); err != nil {
		return err
	}

	createdUser, err := h.sm.CreateUser(c.Context(), &services.CreateUserParams{
		Email:    userBody.Email,
		Name:     userBody.Name,
		Password: userBody.Password,
		Roles:    userBody.Roles,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

type UpdateUserBody struct {
	Status   string   `json:"status"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

func (h *Handlers) updateUser(c *fiber.Ctx) error {
	userId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	userBody := new(UpdateUserBody)

	if err := c.BodyParser(userBody); err != nil {
		return err
	}

	updatedUser, err := h.sm.UpdateUser(c.Context(), &services.UpdateUserParams{
		ID:       userId,
		Status:   userBody.Status,
		Email:    userBody.Email,
		Name:     userBody.Name,
		Password: userBody.Password,
		Roles:    userBody.Roles,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}
