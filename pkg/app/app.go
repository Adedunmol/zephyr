package app

import (
	"fmt"
	"net/http"

	"github.com/Adedunmol/zephyr/pkg/database"
	"github.com/Adedunmol/zephyr/pkg/helpers"
	"github.com/Adedunmol/zephyr/pkg/routes"
)

const PORT = 5001

func init() {
	err := helpers.LoadConfig("../..")

	database.InitDB()

	if err != nil {
		helpers.Error.Fatal("Error loading .env file", err)
	}
}

func Run() {

	addr := fmt.Sprintf(":%d", PORT)

	m := routes.SetupRoutes()

	helpers.Info.Printf("Server listening on: %s", addr)
	http.ListenAndServe(addr, m)
}
