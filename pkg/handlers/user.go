package handlers

import (
	"context"
	"errors"
	"fmt"
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
	Password  string `json:"password" validate:"required,min=6"`
	Email     string `json:"email" validate:"required,email"`
}

func (u *CreateUser) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)

	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Struct(u); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			switch err.Tag() {
			case "required":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' cannot be blank", err.Field())
				problems[field] = message
			case "email":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' must be a valid email address", err.Field())
				problems[field] = message
			case "len":
				field := err.Field()
				message := fmt.Sprintf("Field '%s' must be exactly %v characters long", err.Field(), err.Param())
				problems[field] = message
			default:
				field := err.Field()
				message := fmt.Sprintf("Field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
				problems[field] = message
			}
		}
	}

	return problems
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	data, problems, err := helpers.DecodeJSON[*CreateUser](r)

	if err != nil {

		if err == helpers.ErrValidation {
			helpers.RespondWithJSON(w, http.StatusUnprocessableEntity, helpers.APIResponse{Status: "error", Message: "error processing data", Data: problems})
			return
		}

		if err == helpers.ErrDecode {
			helpers.Error.Println(err)
			helpers.RespondWithJSON(w, http.StatusBadRequest, helpers.APIResponse{Status: "error", Message: "request body needed", Data: nil})
			return
		}
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
