package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const PORT = 5000

func main() {
	r := chi.NewRouter()

	addr := fmt.Sprintf(":%d", PORT)

	http.ListenAndServe(addr, r)
}
