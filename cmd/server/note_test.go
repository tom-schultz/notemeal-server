package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
	"notemeal-server/internal/test"
	"testing"
)

func getNote(id string) *internal.Note {
	note, err := database.Db.GetNote(id)

	if err != nil {
		log.Fatal(err)
	}

	return note
}

func noteSetup() *httptest.Server {
	database.DictDb()

	mux := ServeMux()
	ts := httptest.NewServer(mux)

	return ts
}

func TestNoteDeleteNoAuth(t *testing.T) {
	ts := noteSetup()
	defer ts.CloseClientConnections()
	test.UnauthorizedTest("DELETE", ts.URL+"/notes", nil)
}

func TestNoteDelete(t *testing.T) {
	ts := noteSetup()
	defer ts.CloseClientConnections()
	user := "tom"
	noteId := "dogs"
	token := test.SetupAuth(user)

	note := getNote(noteId)

	if note == nil {
		log.Fatal("Note should not be nil!")
	}

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(user, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	note = getNote(noteId)

	if note != nil {
		log.Fatal("Note should be nil!")
	}
}

func TestNoteGetNoAuth(t *testing.T) {
	ts := noteSetup()
	defer ts.CloseClientConnections()
	test.UnauthorizedTest("GET", ts.URL+"/notes", nil)
}

func TestNoteGet(t *testing.T) {
	ts := noteSetup()
	defer ts.CloseClientConnections()
	user := "tom"
	noteId := "dogs"
	token := test.SetupAuth(user)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(user, token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	note := getNote(noteId)
	test.ExpectBody(resp, test.Serialize(note))
}

func TestNotePutNoAuth(t *testing.T) {
	ts := noteSetup()
	defer ts.CloseClientConnections()
	test.UnauthorizedTest("PUT", ts.URL+"/notes", nil)
}
