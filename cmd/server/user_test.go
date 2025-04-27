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

func getUser(id string) *internal.User {
	user, err := database.Db.GetUser(id)

	if err != nil {
		log.Fatal(err)
	}

	return user
}

func getUserUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/user/%s", baseUrl, id)
}

func TestUserDeleteNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	url := getUserUrl("tom", ts.URL)
	test.UnauthorizedTest("DELETE", url, nil)
}

func TestUserDeleteDifferentPrincipal(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	deleteUserId := "mot"
	url := getUserUrl(deleteUserId, ts.URL)
	token := test.SetupAuth(userId)

	req := test.NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(userId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestUserDelete(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	principalId := "tom"
	objId := "tom"
	url := getUserUrl(objId, ts.URL)
	token := test.SetupAuth(principalId)

	deletedUser := getUser(objId)

	if deletedUser == nil {
		log.Fatal("User to delete does not exist!")
	}

	req := test.NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(principalId, token)
	test.SendReq(req)
	deletedUser = getUser(objId)
	test.ExpectEqual(deletedUser, nil)
}

func TestUserGetNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	url := getUserUrl("tom", ts.URL)
	test.UnauthorizedTest("GET", url, nil)
}

func TestUserGetDifferentPrincipal(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	principalId := "tom"
	objId := "mot"
	url := getUserUrl(objId, ts.URL)
	token := test.SetupAuth(principalId)

	req := test.NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(principalId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestUserGet(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getUserUrl(userId, ts.URL)
	token := test.SetupAuth(userId)

	req := test.NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(userId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	user := getUser(userId)
	test.ExpectBody(resp, test.Serialize(user))
}

func TestUserPutNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	url := getUserUrl("tom", ts.URL)
	test.UnauthorizedTest("PUT", url, nil)
}

func TestUserPutDifferentPrincipal(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	principalId := "tom"
	objId := "mot"
	url := getUserUrl(objId, ts.URL)
	token := test.SetupAuth(principalId)

	req := test.NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(principalId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestUserPut(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getUserUrl(userId, ts.URL)
	token := test.SetupAuth(userId)

	putUser := &internal.User{Id: userId, Email: "new@email.com"}
	req := test.NewReq(http.MethodPut, url, test.Serialize(putUser))
	req.SetBasicAuth(userId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	dbUser := getUser(userId)
	test.ExpectEqual(*putUser, *dbUser)
}

func TestUserPutIdChanged(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	userId := "tom"
	url := getUserUrl(userId, ts.URL)
	token := test.SetupAuth(userId)

	putUser := &internal.User{Id: "malicious", Email: "new@email.com"}
	req := test.NewReq(http.MethodPut, url, test.Serialize(putUser))
	req.SetBasicAuth(userId, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	dbUser := getUser(userId)
	test.ExpectEqual(userId, dbUser.Id)
}
