package main

import (
	"fmt"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/test"
	"testing"
	"time"
)

func buildAdminCodeUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/admin/user/%s/code", baseUrl, id)
}

func TestPutCodeAdminNotAdmin(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	principalId := "tom"
	objId := "tom"
	url := buildAdminCodeUrl(objId, ts.URL)
	token := test.SetupAuth(principalId, m)

	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestCodeAdminPutNew(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	principalId := "admin"
	userId := "mot"
	token := test.SetupAuth(principalId, m)

	prePutCode := getCode(userId, m)
	test.ExpectEqual(prePutCode, nil)

	url := buildAdminCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCode := getCode(userId, m)
	test.ExpectNotEqual(postPutCode, nil)
}

func TestCodeAdminPutUpdate(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	principalId := "admin"
	userId := "tom"
	token := test.SetupAuth(principalId, m)

	createCode(userId, m)
	prePutCode := getCode(userId, m)
	test.ExpectNotEqual(prePutCode, nil)
	prePutExp := prePutCode.Expiration
	// Sometimes the test runs so fast that the times are the same...
	time.Sleep(2 * time.Nanosecond)

	url := buildAdminCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCodeExp := getCode(userId, m).Expiration
	test.ExpectCompareGreater(postPutCodeExp, prePutExp)

	clientCode := internal.ClientCode{}
	test.GetBodyData(resp, &clientCode)

	dbCode := getCode(userId, m)
	test.ExpectNotEqual(dbCode, nil)

	valid := internal.CompareHashAndString(dbCode.Hash, clientCode.Code)
	test.ExpectEqual(valid, true)
}
