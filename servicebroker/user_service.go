package servicebroker

import (
	"context"
	uuid "github.com/jackc/pgx-gofrs-uuid"
	"time"
)

type UserDTO struct {
	ID        *uuid.UUID `json:"id"`
	Status    string     `json:"status"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Roles     string     `json:"roles"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func validateUserName(name string) error {
	return nil
}

type CreateUserParams struct {
	Email    string   `validate:"required,gte=3,lte=80"`
	Name     string   `validate:"required,email"`
	Password string   `validate:"required,gte=6,lte=20"`
	Roles    []string `validate:"required"`
}

func (s *ServiceBroker) CreateUser(ctx context.Context, params *CreateUserParams) (*UserDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	// valid..

	return nil, nil
}

type UpdateUserParams struct {
	ID       *uuid.UUID `validate:"required"`
	Email    string     `validate:"required,gte=3,lte=80"`
	status   string     `validate:"required,custom_status"`
	Name     string     `validate:"required,email"`
	Password string     `validate:"required,gte=6,lte=20"`
	Roles    []string   `validate:"required"`
}

func (s *ServiceBroker) UpdateUser(ctx context.Context, params *UpdateUserParams) (*UserDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	// valid...

	return nil, nil
}

type FetchUsersParams struct {
	statusOptions []string `validate:"custom_status"`
}

type FetchUsersResponse struct {
	totalCount int
	items      []*UserDTO
}

func (s *ServiceBroker) FetchUsers(ctx context.Context, params *FetchUsersParams) (*FetchUsersResponse, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	// get users...

	return nil, nil
}

func (s *ServiceBroker) FetchUserById(ctx context.Context, userID *uuid.UUID) (*UserDTO, error) {

	// user, err := s.repo.GetUserByID(ctx, userID)
	// basically this, but needs to validate for "emptyRows"

	return nil, nil
}

func (s *ServiceBroker) FetchUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	// user, err := s.repo.GetUserByEmail(ctx, email)
	// basically this, but needs to validate for "emptyRows"

	return nil, nil
}
