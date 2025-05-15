package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
	"notemeal-server/internal/test"
	"testing"
	"time"
)

func createCode(userId string, model model.Model) string {
	code, err := model.CreateOrUpdateCode(userId)

	if err != nil {
		log.Fatal(err)
	}

	return code
}

func getCode(id string, m model.Model) *internal.Code {
	code, err := m.GetCode(id)

	if err != nil {
		log.Fatal(err)
	}

	return code
}

func getCodeUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/user/%s/code", baseUrl, id)
}

func TestCodePutNoAuth(t *testing.T) {
	ts, _ := test.Server()
	defer ts.Close()
	url := getCodeUrl("tom", ts.URL)
	test.UnauthorizedTest(http.MethodPut, url, nil)
}

func TestPutCodeDifferentPrincipal(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	principalId := "tom"
	objId := "mot"
	url := getCodeUrl(objId, ts.URL)
	token := test.SetupAuth(principalId, m)

	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestCodePutNew(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	userId := "tom"
	token := test.SetupAuth(userId, m)

	prePutCode := getCode(userId, m)
	test.ExpectEqual(prePutCode, nil)

	url := getCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCode := getCode(userId, m)
	test.ExpectNotEqual(postPutCode, nil)
}

func TestCodePutUpdate(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	userId := "tom"
	token := test.SetupAuth(userId, m)

	createCode(userId, m)
	prePutCode := getCode(userId, m)
	test.ExpectNotEqual(prePutCode, nil)
	prePutExp := prePutCode.Expiration
	// Sometimes the test runs so fast that the times are the same...
	time.Sleep(2 * time.Nanosecond)

	url := getCodeUrl(userId, ts.URL)
	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	postPutCodeExp := getCode(userId, m).Expiration
	test.ExpectCompareGreater(postPutCodeExp, prePutExp)
}
