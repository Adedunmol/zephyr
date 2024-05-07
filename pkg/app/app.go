package app

import (
	"fmt"
	"net/http"

	_ "github.com/Adedunmol/zephyr/pkg/database"
	"github.com/Adedunmol/zephyr/pkg/helpers"
	"github.com/go-chi/chi/v5"
)

const PORT = 5000

func init() {
	err := helpers.LoadConfig("../..")
	if err != nil {
		helpers.Error.Fatal("Error loading .env file", err)
	}
}

func Run() {
	r := chi.NewRouter()

	addr := fmt.Sprintf(":%d", PORT)

	http.ListenAndServe(addr, r)
}
