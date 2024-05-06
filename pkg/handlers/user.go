package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/Adedunmol/zephyr/pkg/database"
	"github.com/Adedunmol/zephyr/pkg/helpers"
	"github.com/Adedunmol/zephyr/pkg/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUser struct {
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

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	data, problems, err := helpers.DecodeJSON[*CreateUser](r)

	if err != nil && err == helpers.ErrValidation {
		helpers.RespondWithJSON(w, http.StatusUnprocessableEntity, problems)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	if err != nil {
		helpers.Info.Println("could not hash password", err)
		helpers.RespondWithJSON(w, http.StatusInternalServerError, helpers.APIResponse{Message: "Unable to hash password", Data: nil, Status: "error"})
		return
	}

	user := models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Username:  data.Username,
		Password:  string(hashedPassword),
		Email:     data.Email,
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			helpers.Error.Println(result.Error)
			helpers.RespondWithJSON(w, http.StatusConflict, helpers.APIResponse{Message: "Duplicate field sent", Data: nil, Status: "error"})
			return
		}
	}

	helpers.RespondWithJSON(w, http.StatusCreated, helpers.APIResponse{Message: "", Data: user, Status: "success"})
}
