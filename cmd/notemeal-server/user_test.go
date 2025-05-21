package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
	"testing"
)

func buildUserUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/user/%s", baseUrl, id)
}

func getUser(id string, m model.Model) *internal.User {
	user, err := m.GetUser(id)

	if err != nil {
		log.Fatal(err)
	}

	return user
}

func TestUserDeleteNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	url := buildUserUrl("tom", ts.URL)
	UnauthorizedTest("DELETE", url, nil)
}

func TestUserDeleteDifferentPrincipal(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	userId := "tom"
	deleteUserId := "mot"
	url := buildUserUrl(deleteUserId, ts.URL)
	token := SetupAuth(userId, m)

	req := NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestUserDelete(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	principalId := "tom"
	objId := "tom"
	url := buildUserUrl(objId, ts.URL)
	token := SetupAuth(principalId, m)

	deletedUser := getUser(objId, m)

	if deletedUser == nil {
		log.Fatal("User to delete does not exist!")
	}

	req := NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	SendReq(req)
	deletedUser = getUser(objId, m)
	ExpectEqual(deletedUser, nil)
}

func TestUserGetNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	url := buildUserUrl("tom", ts.URL)
	UnauthorizedTest("GET", url, nil)
}

func TestUserGetDifferentPrincipal(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	principalId := "tom"
	objId := "mot"
	url := buildUserUrl(objId, ts.URL)
	token := SetupAuth(principalId, m)

	req := NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestUserGet(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	userId := "tom"
	url := buildUserUrl(userId, ts.URL)
	token := SetupAuth(userId, m)

	req := NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	user := getUser(userId, m)
	ExpectBody(resp, Serialize(user))
}

func TestUserPutNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	url := buildUserUrl("tom", ts.URL)
	UnauthorizedTest("PUT", url, nil)
}

func TestUserPutDifferentPrincipal(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	principalId := "tom"
	objId := "mot"
	url := buildUserUrl(objId, ts.URL)
	token := SetupAuth(principalId, m)

	req := NewReq(http.MethodPut, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusUnauthorized)
}

func TestUserPut(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	userId := "tom"
	url := buildUserUrl(userId, ts.URL)
	token := SetupAuth(userId, m)

	putUser := &internal.User{Id: userId, Email: "new@email.com"}
	req := NewReq(http.MethodPut, url, Serialize(putUser))
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	dbUser := getUser(userId, m)
	ExpectEqual(*putUser, *dbUser)
}

func TestUserPutIdChanged(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	userId := "tom"
	url := buildUserUrl(userId, ts.URL)
	token := SetupAuth(userId, m)

	putUser := &internal.User{Id: "malicious", Email: "new@email.com"}
	req := NewReq(http.MethodPut, url, Serialize(putUser))
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	dbUser := getUser(userId, m)
	ExpectEqual(userId, dbUser.Id)
}
