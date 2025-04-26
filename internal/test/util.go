package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"notemeal-server/internal/database"
)

func ExpectBody(resp *http.Response, body []byte) {
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if bytes.Compare(respBody, body) != 0 {
		log.Fatal("Response body did not match expected body!")
	}
}

func ExpectStatusCode(resp *http.Response, code int) {
	if resp.StatusCode != code {
		log.Fatalf("Expected %d, got %d!", resp.StatusCode, code)
	}
}

func NewReq(method string, url string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))

	if err != nil {
		log.Fatal(err)
	}

	return req
}

func SendReq(req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func SetupAuth(user string) string {
	code, err := database.Db.CreateOrUpdateCode(user)

	if err != nil {
		log.Fatal(err)
	}

	token, err := database.Db.CreateToken(user, code)

	if err != nil {
		log.Fatal(err)
	}

	return token
}

func Serialize(v any) []byte {
	data, err := json.Marshal(v)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func UnauthorizedTest(method string, url string, body []byte) {
	req := NewReq("GET", url, body)
	resp := SendReq(req)

	if resp.StatusCode != http.StatusUnauthorized {
		log.Fatal("Unauthorized request didn't give 401!")
	}
}
