package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Validator interface {
	Valid(ctx context.Context) (problems map[string]string)
}

var ErrValidation error

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.MarshalIndent(payload, "", "   ")
	if err != nil {
		log.Println("Failed to marshal JSON response: ", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

	return
}

func DecodeJSON[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	if problems := v.Valid(r.Context()); len(problems) != 0 {
		ErrValidation = errors.New(fmt.Sprintf("invalid %T: %d problem(s)", v, len(problems)))
		return v, problems, ErrValidation
	}

	return v, nil, nil
}
