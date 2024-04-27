package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mux *http.ServeMux

const checkMark = "\u2713"
const ballotX = "\u2717"

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *CreateUser) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)

	if r.Username == "" {
		problems["username"] = "username required"
	}
	if r.Password == "" {
		problems["password"] = "password required"
	}

	return problems
}

func mockServer() (*httptest.Server, *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/respond-with-json", func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Status string
			Data   string
		}{
			Status: "success",
			Data:   "Hello",
		}
		RespondWithJSON(w, http.StatusOK, payload)
	})

	mux.HandleFunc("/decode-json", func(w http.ResponseWriter, r *http.Request) {
		data, _, err := DecodeJSON[*CreateUser](r)

		if err != nil && err == ErrValidation {
			RespondWithJSON(w, http.StatusUnprocessableEntity, err)
		}

		RespondWithJSON(w, http.StatusOK, data)
	})

	return httptest.NewServer(mux), mux
}

func TestRespondWithJSON(t *testing.T) {
	statusCode := http.StatusOK

	server, mux := mockServer()
	defer server.Close()

	t.Log("Given the need to test responding with JSON.")
	{
		t.Log("\tWhen checking for status code.")
		{
			req, err := http.NewRequest(http.MethodGet, server.URL+"/respond-with-json", nil)
			// res, err := http.Get(server.URL)

			if err != nil {
				t.Fatal("\t\tShould be able to create the GET call", ballotX, err)
			}
			t.Log("\t\tShould be able to create the GET call.", checkMark)

			rw := httptest.NewRecorder()
			mux.ServeHTTP(rw, req)

			if rw.Code != statusCode {
				t.Errorf("\t\tShould receive a %d status code, but got %v. %v", statusCode, rw.Code, ballotX)
			}

			t.Logf("\t\tShould receive a %d status code. %v", statusCode, checkMark)

			payload := struct {
				Status string
				Data   string
			}{}

			if err := json.NewDecoder(rw.Body).Decode(&payload); err != nil {
				t.Fatal("\t\t\tShould decode response.", ballotX)
			}
			t.Log("\t\t\tShould decode response.", checkMark)

			if payload.Status == "success" {
				t.Log("\t\t\t\tShould have status.", checkMark)
			} else {
				t.Log("\t\t\t\tShould have status.", ballotX, payload.Status)
			}

			if payload.Data == "Hello" {
				t.Log("\t\t\t\tShould have data.", checkMark)
			} else {
				t.Log("\t\t\t\tShould have data.", ballotX, payload.Data)
			}
		}
	}
}

func TestDecodeJSON(t *testing.T) {
	statusCode := http.StatusOK

	server, mux := mockServer()
	defer server.Close()

	t.Log("Given the need to test decoding JSON.")
	{
		t.Log("When checking for a passing validation")
		{
			reqBody := `
			{
				"username": "Adedunmola", "password": "password123"
			}
			`

			req, err := http.NewRequest(http.MethodPost, server.URL+"/decode-json", strings.NewReader(reqBody))

			if err != nil {
				t.Fatal("\t\tShould be able to create the POST call", ballotX, err)
			}
			t.Log("\t\tShould be able to create the POST call.", checkMark)

			rw := httptest.NewRecorder()
			mux.ServeHTTP(rw, req)

			if rw.Code != statusCode {
				t.Log(rw.Body)
				t.Errorf("\t\tShould receive a %d status code, but got %v. %v", statusCode, rw.Code, ballotX)
			}

			t.Logf("\t\tShould receive a %d status code. %v", statusCode, checkMark)
		}

		statusCode = http.StatusUnprocessableEntity
		t.Log("When checking for a failing validation")
		{
			reqBody := `
			{
				"username": "Adedunmola", "password": ""
			}
			`

			req, err := http.NewRequest(http.MethodPost, server.URL+"/decode-json", strings.NewReader(reqBody))

			if err != nil {
				t.Fatal("\t\tShould be able to create the POST call", ballotX, err)
			}
			t.Log("\t\tShould be able to create the POST call.", checkMark)

			rw := httptest.NewRecorder()
			mux.ServeHTTP(rw, req)

			if rw.Code != statusCode {
				t.Log(rw.Body)
				t.Errorf("\t\tShould receive a %d status code, but got %v. %v", statusCode, rw.Code, ballotX)
			}

			t.Logf("\t\tShould receive a %d status code. %v", statusCode, checkMark)
		}
	}
}
