package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
	"testing"
)

func buildTokenUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/user/%s/token", baseUrl, id)
}

func getToken(id string, m model.Model) *internal.Token {
	token, err := m.GetToken(id)

	if err != nil {
		log.Fatal(err)
	}

	return token
}

func TestTokenPost(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	userId := "tom"
	url := buildTokenUrl(userId, ts.URL)
	code := createCode(userId, m)
	clientCode := internal.ClientCode{UserId: userId, Code: code}
	reqBody := Serialize(clientCode)

	req := NewReq(http.MethodPost, url, reqBody)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	var clientToken internal.ClientToken
	GetBodyData(resp, &clientToken)
	dbToken := getToken(clientToken.Id, m)
	ExpectNotEqual(dbToken, nil)
	valid := internal.CompareHashAndString(dbToken.Hash, clientToken.Token)
	ExpectEqual(valid, true)
}

func TestTokenPostFakeCode(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	userId := "tom"
	url := buildTokenUrl(userId, ts.URL)
	code := "mumbojumbo"
	clientCode := internal.ClientCode{UserId: userId, Code: code}
	reqBody := Serialize(clientCode)

	req := NewReq(http.MethodPost, url, reqBody)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestTokenPostWrongCode(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	userId := "tom"
	url := buildTokenUrl(userId, ts.URL)
	createCode(userId, m)
	wrongCode := "mumbojumbo"
	clientCode := internal.ClientCode{UserId: userId, Code: wrongCode}
	reqBody := Serialize(clientCode)

	req := NewReq(http.MethodPost, url, reqBody)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestTokenPostExpiredCode(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	userId := "expired-code"
	url := buildTokenUrl(userId, ts.URL)
	code := "expired"
	clientCode := internal.ClientCode{UserId: userId, Code: code}
	reqBody := Serialize(clientCode)

	req := NewReq(http.MethodPost, url, reqBody)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}
