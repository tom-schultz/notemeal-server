package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/model"
	"notemeal-server/internal/test"
	"testing"
)

func listNotes(user string, m model.Model) []byte {
	notes, err := m.ListLastModified(user)

	if err != nil {
		log.Fatal(err)
	}

	return test.Serialize(notes)
}

func TestNotesGetNoAuth(t *testing.T) {
	ts, _ := test.Server()
	defer ts.Close()
	test.UnauthorizedTest("GET", ts.URL+"/notes", nil)
}

func TestNotesGet(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	user := "tom"
	token := test.SetupAuth(user, m)

	req := test.NewReq("GET", ts.URL+"/notes", nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	notes := listNotes(user, m)
	test.ExpectBody(resp, notes)
}
