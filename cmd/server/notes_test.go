package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"notemeal-server/internal/database"
	"notemeal-server/internal/handler"
	"notemeal-server/internal/test"
	"testing"
)

func listNotes(user string) []byte {
	notes, err := database.Db.ListLastModified(user)

	if err != nil {
		log.Fatal(err)
	}

	return test.Serialize(notes)
}

func notesSetup() *httptest.Server {
	database.DictDb()
	ts := httptest.NewServer(http.HandlerFunc(handler.GetNotes))
	return ts
}

func TestNotesGetNoAuth(t *testing.T) {
	ts := notesSetup()
	defer ts.CloseClientConnections()
	test.UnauthorizedTest("PUT", ts.URL+"/notes", nil)
}

func TestNotesGet(t *testing.T) {
	ts := notesSetup()
	defer ts.CloseClientConnections()
	user := "tom"
	token := test.SetupAuth(user)

	req := test.NewReq("PUT", ts.URL+"/notes", nil)
	req.SetBasicAuth(user, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	notes := listNotes(user)
	test.ExpectBody(resp, notes)
}
