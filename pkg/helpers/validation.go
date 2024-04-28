package helpers

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type Validate struct {
	Schema interface{}
}

func (d *Validate) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)

	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(d.Schema); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			message := err.Tag()

			problems[field] = message
		}
	}

	return problems
}
