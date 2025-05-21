package main

import (
	"fmt"
	"net/http"
	"testing"
)

func buildHealthUrl(baseUrl string) string {
	return fmt.Sprintf("%s/health", baseUrl)
}

func TestGetHealthNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	url := buildHealthUrl(ts.URL)

	req := NewReq(http.MethodGet, url, nil)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)
}

func TestGetHealthAuth(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	principalId := "tom"
	url := buildHealthUrl(ts.URL)
	token := SetupAuth(principalId, m)

	req := NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)
}
