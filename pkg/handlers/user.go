package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Adedunmol/zephyr/pkg/database"
	"github.com/Adedunmol/zephyr/pkg/helpers"
	"github.com/Adedunmol/zephyr/pkg/models"
	"github.com/Adedunmol/zephyr/pkg/schema"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	data, problems, err := helpers.DecodeJSON[*schema.CreateUser](r)

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

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token      string        `json:"token"`
		Expiration time.Duration `json:"expiration"`
	}

	data, problems, err := helpers.DecodeJSON[*schema.CreateUser](r)

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

	var foundUser models.User

	result := database.DB.Where(models.User{Email: data.Email}).First(&foundUser)

	if result.Error != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, helpers.APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(data.Password))

	if err != nil {
		helpers.RespondWithJSON(w, http.StatusUnauthorized, helpers.APIResponse{Message: "Invalid credentials", Data: nil, Status: "error"})
		return
	}

	accessToken, err := helpers.GenerateToken(foundUser.Username, helpers.ACCESS_TOKEN_EXPIRATION)

	if err != nil {
		helpers.Error.Println(err)
		helpers.RespondWithJSON(w, http.StatusInternalServerError, helpers.APIResponse{Message: "Unable to generate token", Data: nil, Status: "error"})
		return
	}

	refreshToken, err := helpers.GenerateToken(foundUser.Username, helpers.REFRESH_TOKEN_EXPIRATION)

	if err != nil {
		helpers.Error.Println(err)
		helpers.RespondWithJSON(w, http.StatusInternalServerError, helpers.APIResponse{Message: "Unable to generate token", Data: nil, Status: "error"})
		return
	}

	cookie := http.Cookie{
		Name:  "token",
		Value: refreshToken,
		// Expires:  time.Now().Add(util.REFRESH_TOKEN_EXPIRATION),
		HttpOnly: true,
		MaxAge:   1 * 60 * 60,
	}

	result = database.DB.Model(&foundUser).UpdateColumn("RefreshToken", refreshToken)

	if result.Error != nil {
		helpers.Error.Println(result.Error)
		helpers.RespondWithJSON(w, http.StatusInternalServerError, helpers.APIResponse{Message: "unable to update refresh token", Data: nil, Status: "error"})
		return
	}

	res := Response{Token: accessToken, Expiration: time.Duration(helpers.ACCESS_TOKEN_EXPIRATION.Seconds())}

	http.SetCookie(w, &cookie)
	helpers.RespondWithJSON(w, http.StatusOK, helpers.APIResponse{Message: "", Data: res, Status: "success"})
}
