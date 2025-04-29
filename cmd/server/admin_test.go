package main

import (
	"fmt"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
	"notemeal-server/internal/test"
	"testing"
	"time"
)

func getAdminCodeUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/admin/user/%s/code", baseUrl, id)
}

func TestPutCodeAdminNotAdmin(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	principalId := "tom"
	objId := "tom"
	url := getAdminCodeUrl(objId, ts.URL)
	token := test.SetupAuth(principalId)

	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestCodeAdminPutNew(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	principalId := "admin"
	userId := "tom"
	token := test.SetupAuth(principalId)

	prePutCode := getCode(userId)
	test.ExpectEqual(prePutCode, nil)

	url := getAdminCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCode := getCode(userId)
	test.ExpectNotEqual(postPutCode, nil)
}

func TestCodeAdminPutUpdate(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	principalId := "admin"
	userId := "tom"
	token := test.SetupAuth(principalId)

	createCode(userId)
	prePutCode := getCode(userId)
	test.ExpectNotEqual(prePutCode, nil)
	prePutExp := prePutCode.Expiration
	// Sometimes the test runs so fast that the times are the same...
	time.Sleep(2 * time.Nanosecond)

	url := getAdminCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCodeExp := getCode(userId).Expiration
	test.ExpectCompareGreater(postPutCodeExp, prePutExp)

	var bodyData map[string]string
	test.GetBodyData(resp, &bodyData)

	dbCode := getCode(userId)
	test.ExpectNotEqual(dbCode, nil)

	err := database.CompareHashAndString(dbCode.Hash, bodyData[internal.CodeJsonKey])
	test.ExpectEqual(err, nil)
}
