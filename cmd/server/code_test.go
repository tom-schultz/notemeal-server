package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
	"notemeal-server/internal/test"
	"testing"
	"time"
)

func createCode(userId string) string {
	code, err := database.Db.CreateOrUpdateCode(userId)

	if err != nil {
		log.Fatal(err)
	}

	return code
}

func getCode(id string) *internal.Code {
	code, err := database.Db.GetCode(id)

	if err != nil {
		log.Fatal(err)
	}

	return code
}

func getCodeUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/user/%s/code", baseUrl, id)
}

func TestCodePutNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	url := getCodeUrl("tom", ts.URL)
	test.UnauthorizedTest(http.MethodPut, url, nil)
}

func TestPutCodeDifferentPrincipal(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	principalId := "tom"
	objId := "mot"
	url := getCodeUrl(objId, ts.URL)
	token := test.SetupAuth(principalId)

	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(principalId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestCodePutNew(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	token := test.SetupAuth(userId)

	prePutCode := getCode(userId)
	test.ExpectEqual(prePutCode, nil)

	url := getCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(userId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCode := getCode(userId)
	test.ExpectNotEqual(postPutCode, nil)
}

func TestCodePutUpdate(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	token := test.SetupAuth(userId)

	createCode(userId)
	prePutCode := getCode(userId)
	test.ExpectNotEqual(prePutCode, nil)
	prePutExp := prePutCode.Expiration
	// Sometimes the test runs so fast that the times are the same...
	time.Sleep(2 * time.Nanosecond)

	url := getCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(userId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCodeExp := getCode(userId).Expiration
	test.ExpectCompareGreater(postPutCodeExp, prePutExp)
}
