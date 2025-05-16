package main

import (
	"fmt"
	"net/http"
	"notemeal-server/internal"
	"testing"
	"time"
)

func buildAdminCodeUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/admin/user/%s/code", baseUrl, id)
}

func TestPutCodeAdminNotAdmin(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	principalId := "tom"
	objId := "tom"
	url := buildAdminCodeUrl(objId, ts.URL)
	token := SetupAuth(principalId, m)

	req := NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestCodeAdminPutNew(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	principalId := "admin"
	userId := "mot"
	token := SetupAuth(principalId, m)

	prePutCode := getCode(userId, m)
	ExpectEqual(prePutCode, nil)

	url := buildAdminCodeUrl(userId, ts.URL)
	req := NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	postPutCode := getCode(userId, m)
	ExpectNotEqual(postPutCode, nil)
}

func TestCodeAdminPutUpdate(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	principalId := "admin"
	userId := "tom"
	token := SetupAuth(principalId, m)

	createCode(userId, m)
	prePutCode := getCode(userId, m)
	ExpectNotEqual(prePutCode, nil)
	prePutExp := prePutCode.Expiration
	// Sometimes the test runs so fast that the times are the same...
	time.Sleep(2 * time.Nanosecond)

	url := buildAdminCodeUrl(userId, ts.URL)
	req := NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	postPutCodeExp := getCode(userId, m).Expiration
	ExpectCompareGreater(postPutCodeExp, prePutExp)

	clientCode := internal.ClientCode{}
	GetBodyData(resp, &clientCode)

	dbCode := getCode(userId, m)
	ExpectNotEqual(dbCode, nil)

	valid := internal.CompareHashAndString(dbCode.Hash, clientCode.Code)
	ExpectEqual(valid, true)
}
