package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
	"notemeal-server/internal/test"
	"testing"
)

func getToken(id string) *internal.Token {
	token, err := database.Db.GetToken(id)

	if err != nil {
		log.Fatal(err)
	}

	return token
}

func getTokenUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/user/%s/token", baseUrl, id)
}

func TestTokenPost(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getTokenUrl(userId, ts.URL)
	code := createCode(userId)
	codeMap := map[string]string{internal.CodeJsonKey: code}
	reqBody := test.Serialize(codeMap)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	respData := make(map[string]string)
	test.GetBodyData(resp, &respData)
	respToken := respData[internal.TokenJsonKey]
	dbToken := getToken(respToken)
	test.ExpectNotEqual(dbToken, nil)
	test.ExpectEqual(database.HashString(respToken), dbToken.TokenHash)
}

func TestTokenPostFakeCode(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getTokenUrl(userId, ts.URL)
	code := "mumbojumbo"
	codeMap := map[string]string{internal.CodeJsonKey: code}
	reqBody := test.Serialize(codeMap)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestTokenPostWrongCode(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getTokenUrl(userId, ts.URL)
	createCode(userId)
	wrongCode := "mumbojumbo"
	codeMap := map[string]string{internal.CodeJsonKey: wrongCode}
	reqBody := test.Serialize(codeMap)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestTokenPostExpiredCode(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "expired-code"
	url := getTokenUrl(userId, ts.URL)
	code := "expired"
	codeMap := map[string]string{internal.CodeJsonKey: code}
	reqBody := test.Serialize(codeMap)

	req := test.NewReq(http.MethodPost, url, reqBody)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}
