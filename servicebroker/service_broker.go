package servicebroker

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/hoffax/prodrest/repository"
)

type ServiceBroker struct {
	repo     *repository.PgRepository
	validate *validator.Validate
}

func NewServiceBroker(repo *repository.PgRepository) (*ServiceBroker, error) {
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

	return &ServiceBroker{
		repo:     repo,
		validate: validate,
	}, nil
}
