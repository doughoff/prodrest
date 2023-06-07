package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"
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
	if dbuuid == nil {
		return nil, fmt.Errorf("dbuuid cannot be nil")
	}

	uuidValue, err := dbuuid.UUIDValue()
	if err != nil {
		return nil, err
	}

	result := uuid.UUID(uuidValue.Bytes)
	return &result, nil
}

func (s *ServiceManager) comparePgxUUID(uuid1 *pgxuuid.UUID, uuid2 *pgxuuid.UUID) (bool, error) {
	uuid1Value, err := s.parseUUID(uuid1)
	if err != nil {
		return false, err
	}

	uuid2Value, err := s.parseUUID(uuid2)
	if err != nil {
		return false, err
	}

	return bytes.Equal(uuid1Value.Bytes(), uuid2Value.Bytes()), nil
}
