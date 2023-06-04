package services

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/hoffax/prodrest/constants"
	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type UserDTO struct {
	ID        *uuid.UUID `json:"id"`
	Status    string     `json:"status"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Roles     []string   `json:"roles"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (s *ServiceManager) toUserDTO(user *repository.User) *UserDTO {
	userId, err := s.parseUUID(user.ID)
	if err != nil {
		userId = nil
	}

	return &UserDTO{
		ID:        userId,
		Status:    user.Status,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type CreateUserParams struct {
	Email    string   `validate:"required,email"`
	Name     string   `validate:"required,gte=3,lte=80"`
	Password string   `validate:"required,gte=6,lte=20"`
	Roles    []string `validate:"required"`
}

func (s *ServiceManager) CreateUser(ctx context.Context, params *CreateUserParams) (*UserDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.GetUserByEmail(ctx, params.Email)
	if err == nil {
		return nil, constants.NewUniqueConstrainError("email")
	} else {
		if err != pgx.ErrNoRows {
			return nil, err
		}
	}

	user, err := s.repo.CreateUser(ctx, &repository.NewUserParams{
		Email:    params.Email,
		Name:     params.Name,
		Password: params.Password,
		Roles:    params.Roles,
	})
	if err != nil {
		return nil, err
	}

	return s.toUserDTO(user), nil
}

type UpdateUserParams struct {
	ID       *pgxuuid.UUID `validate:"required"`
	Status   string        `validate:"required,custom_status"`
	Email    string        `validate:"required,email"`
	Name     string        `validate:"required,gte=3,lte=80"`
	Password string        `validate:"required,gte=6,lte=20"`
	Roles    []string      `validate:"required"`
}

func (s *ServiceManager) UpdateUser(ctx context.Context, params *UpdateUserParams) (*UserDTO, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.GetUserByID(ctx, params.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, constants.NewNotFoundError()
		}
		return nil, err
	}

	userWithEmail, err := s.repo.GetUserByEmail(ctx, params.Email)
	if err == nil {
		if userWithEmail.ID != params.ID {
			return nil, constants.NewUniqueConstrainError("email")
		}
	} else {
		if err != pgx.ErrNoRows {
			return nil, err
		}
	}

	updatedUser, err := s.repo.UpdateUser(ctx, &repository.UpdateUserParams{
		ID:       params.ID,
		Status:   params.Status,
		Email:    params.Email,
		Name:     params.Name,
		Password: params.Password,
		Roles:    params.Roles,
	})
	if err != nil {
		return nil, err
	}

	return s.toUserDTO(updatedUser), nil
}

type FetchUsersParams struct {
	StatusOptions []string `validate:"dive,custom_status"`
}

type FetchUsersResponse struct {
	TotalCount int        `json:"totalCount"`
	Items      []*UserDTO `json:"items"`
}

func (s *ServiceManager) FetchUsers(ctx context.Context, params *FetchUsersParams) (*FetchUsersResponse, error) {
	err := s.validate.Struct(params)
	if err != nil {
		return nil, err
	}

	if len(params.StatusOptions) == 0 {
		params.StatusOptions = []string{"ACTIVE", "INACTIVE"}
	}

	users, err := s.repo.GetAllUsers(ctx, &repository.GetAllUsersParams{
		Status: params.StatusOptions,
	})
	if err != nil {
		return nil, err
	}

	usersDTO := make([]*UserDTO, len(users))
	for i := 0; i < len(users); i++ {
		usersDTO[i] = s.toUserDTO(users[i])
	}

	return &FetchUsersResponse{
		TotalCount: len(usersDTO),
		Items:      usersDTO,
	}, nil
}

func (s *ServiceManager) FetchUserById(ctx context.Context, userID *pgxuuid.UUID) (*UserDTO, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.NewNotFoundError()
		}
		return nil, err
	}

	return s.toUserDTO(user), nil
}

func (s *ServiceManager) FetchUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.NewNotFoundError()
		}
		return nil, err
	}

	return s.toUserDTO(user), nil
}

func (s *ServiceManager) CheckEmailAndPassword(ctx context.Context, email string, password string) (bool, *UserDTO, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil, nil
		}
		return false, nil, err
	}

	if user.Status != "ACTIVE" {
		return false, nil, nil
	}

	if user.Password == password {
		return true, s.toUserDTO(user), nil
	}

	return false, nil, nil
}
