package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreateUser struct {
	gorm.Model
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required, min=6"`
	Email     string `json:"email" validate:"required, email"`
}

func (u *CreateUser) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)

	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(u); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			message := err.Tag()

			problems[field] = message
		}
	}

	return problems
}
