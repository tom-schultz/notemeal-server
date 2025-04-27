package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/database"
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

func TestNotesGetNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	test.UnauthorizedTest("GET", ts.URL+"/notes", nil)
}

func TestNotesGet(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	user := "tom"
	token := test.SetupAuth(user)

	req := test.NewReq("GET", ts.URL+"/notes", nil)
	req.SetBasicAuth(user, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	notes := listNotes(user)
	test.ExpectBody(resp, notes)
}
