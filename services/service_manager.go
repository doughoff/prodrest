package services

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/hoffax/prodrest/repository"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
)

type ServiceManager struct {
	repo     *repository.PgRepository
	validate *validator.Validate
}

func NewServiceManager(repo *repository.PgRepository) (*ServiceManager, error) {
	validate := validator.New()
	err := validate.RegisterValidation("custom_status", func(fl validator.FieldLevel) bool {
		value := fl.Field()

		return value.String() != "unknown" && (value.String() == "ACTIVE" || value.String() == "INACTIVE")
	})
	if err != nil {
		return nil, errors.New("could not load custom_status validator")
	}

	err = validate.RegisterValidation("custom_unit", func(fl validator.FieldLevel) bool {
		value := fl.Field()

		return value.String() != "unknown" && (value.String() == "KG" || value.String() == "L" ||
			value.String() == "UN" || value.String() == "OTHER")
	})
	if err != nil {
		return nil, errors.New("could not load custom_unit validator")
	}

	return &ServiceManager{
		repo:     repo,
		validate: validate,
	}, nil
}

func (s *ServiceManager) parseUUID(dbuuid *pgxuuid.UUID) (*uuid.UUID, error) {
	uuidValue, err := dbuuid.UUIDValue()
	if err != nil {
		return nil, err
	}

	result := uuid.UUID(uuidValue.Bytes)
	return &result, nil
}

type UniqueConstraintError struct {
	Message string
}

func (u UniqueConstraintError) Error() string {
	return u.Message
}

func NewUniqueConstrainError(field string) *UniqueConstraintError {
	return &UniqueConstraintError{
		fmt.Sprintf("unique constraint error on: %v", field),
	}
}
