package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
	"notemeal-server/internal/handler"
	"notemeal-server/internal/model"
)

type Comparer[T any] interface {
	Compare(T) int
}

func ExpectBody(resp *http.Response, body []byte) {
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if bytes.Compare(respBody, body) != 0 {
		log.Fatal("Response body did not match expected body!")
	}
}

func ExpectEqual[T comparable](lhs T, rhs T) {
	if lhs != rhs {
		log.Fatal("Expected equal!")
	}
}

func ExpectNotEqual[T comparable](lhs T, rhs T) {
	if lhs == rhs {
		log.Fatal("Expected not equal!")
	}
}

func ExpectCompareGreater[T Comparer[T]](lhs T, rhs T) {
	if lhs.Compare(rhs) <= 0 {
		log.Fatal("Expected not equal!")
	}
}

func ExpectStatusCode(resp *http.Response, code int) {
	if resp.StatusCode != code {
		log.Fatalf("Expected %d, got %d!", code, resp.StatusCode)
	}
}

func GetBodyData[T any](resp *http.Response, respData *T) {
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Could not read response body!")
	}

	err = json.Unmarshal(respBody, respData)

	if err != nil {
		log.Fatal("Could not deserialize response body!")
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

func SetupAuth(user string, model model.Model) internal.ClientToken {
	code, err := model.CreateOrUpdateCode(user)

	if err != nil {
		log.Fatal(err)
	}

	token, err := model.CreateToken(user, code)

	if err != nil {
		log.Fatal(err)
	}

	return *token
}

func Serialize(v any) []byte {
	data, err := json.Marshal(v)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func Server() (*httptest.Server, model.Model) {
	db := database.DictDb()
	m := model.NewModel(db)

	mux := handler.ServeMux(m)
	ts := httptest.NewServer(mux)

	return ts, m
}

func UnauthorizedTest(method string, url string, body []byte) {
	req := NewReq(method, url, body)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}
