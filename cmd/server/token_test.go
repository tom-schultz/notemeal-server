package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
	"notemeal-server/internal/test"
	"testing"
)

func getToken(id string, m model.Model) *internal.Token {
	token, err := m.GetToken(id)

	if err != nil {
		log.Fatal(err)
	}

	return token
}

func getTokenUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/user/%s/token", baseUrl, id)
}

func TestTokenPost(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getTokenUrl(userId, ts.URL)
	code := createCode(userId, m)
	clientCode := internal.ClientCode{UserId: userId, Code: code}
	reqBody := test.Serialize(clientCode)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	var clientToken internal.ClientToken
	test.GetBodyData(resp, &clientToken)
	dbToken := getToken(clientToken.Id, m)
	test.ExpectNotEqual(dbToken, nil)
	err := internal.CompareHashAndString(dbToken.Hash, clientToken.Token)
	test.ExpectEqual(err, nil)
}

func TestTokenPostFakeCode(t *testing.T) {
	ts, _ := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getTokenUrl(userId, ts.URL)
	code := "mumbojumbo"
	clientCode := internal.ClientCode{UserId: userId, Code: code}
	reqBody := test.Serialize(clientCode)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestTokenPostWrongCode(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getTokenUrl(userId, ts.URL)
	createCode(userId, m)
	wrongCode := "mumbojumbo"
	clientCode := internal.ClientCode{UserId: userId, Code: wrongCode}
	reqBody := test.Serialize(clientCode)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestTokenPostExpiredCode(t *testing.T) {
	ts, _ := test.Server()
	defer ts.Close()
	userId := "expired-code"
	url := getTokenUrl(userId, ts.URL)
	code := "expired"
	clientCode := internal.ClientCode{UserId: userId, Code: code}
	reqBody := test.Serialize(clientCode)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}
