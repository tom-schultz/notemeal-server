package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/model"
	"testing"
)

func listNotes(user string, m model.Model) []byte {
	notes, err := m.ListLastModified(user)

	if err != nil {
		log.Fatal(err)
	}

	return Serialize(notes)
}

func TestNotesGetNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	UnauthorizedTest("GET", ts.URL+"/notes", nil)
}

func TestNotesGet(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	user := "tom"
	token := SetupAuth(user, m)

	req := NewReq("GET", ts.URL+"/notes", nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	notes := listNotes(user, m)
	ExpectBody(resp, notes)
}
