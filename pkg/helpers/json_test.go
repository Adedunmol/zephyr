package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Status string
			Data   string
		}{
			Status: "success",
			Data:   "Hello",
		}
		RespondWithJSON(w, http.StatusOK, payload)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestRespondWithJSON(t *testing.T) {
	statusCode := http.StatusOK

	server := mockServer()
	defer server.Close()

	t.Log("Given the need to test responding with JSON.")
	{
		t.Logf("\tWhen checking %s for %d status code.", server.URL, statusCode)
		{
			res, err := http.Get(server.URL)

			if err != nil {
				t.Fatal("\t\tShould be able to make the GET call", ballotX, err)
			}
			t.Log("\t\tShould be able to make the GET call.", checkMark)

			if res.StatusCode != statusCode {
				t.Errorf("\t\tShould receive a %d status code, but got %v. %v", statusCode, res.StatusCode, ballotX)
			}

			t.Logf("\t\tShould receive a %d status code. %v", statusCode, checkMark)

			payload := struct {
				Status string
				Data   string
			}{}

			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
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
