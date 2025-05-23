package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
	"testing"
)

func buildNoteUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/note/%s", baseUrl, id)
}

func getNote(id string, m model.Model) *internal.Note {
	note, err := m.GetNote(id)

	if err != nil {
		log.Fatal(err)
	}

	return note
}

func TestNoteDeleteNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	url := buildNoteUrl("dogs", ts.URL)
	UnauthorizedTest("DELETE", url, nil)
}

func TestNoteDelete(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	token := SetupAuth(user, m)

	note := getNote(noteId, m)

	if note == nil {
		log.Fatal("Test note does not exist!")
	}

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	note = getNote(noteId, m)

	if note != nil {
		log.Fatal("Note should be nil!")
	}
}

func TestNoteGetNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	url := buildNoteUrl("dogs", ts.URL)
	UnauthorizedTest("GET", url, nil)
}

func TestNoteGet(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	token := SetupAuth(user, m)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	note := getNote(noteId, m)
	ExpectBody(resp, Serialize(note))
}

func TestNotePutCreate(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	user := "tom"
	noteId := "wuppers"
	postNote := &internal.Note{Id: noteId, Text: "woof woof", Title: "Puppers", UserId: user, LastModified: 0}
	token := SetupAuth(user, m)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := NewReq(http.MethodPut, url, Serialize(postNote))
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	dbNote := getNote(noteId, m)
	ExpectEqual(*dbNote, *postNote)
}

func TestNotePutUpdate(t *testing.T) {
	ts, m := Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	putNote := &internal.Note{Id: noteId, Text: "woof woof", Title: "Puppers", UserId: user, LastModified: 0}
	token := SetupAuth(user, m)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := NewReq(http.MethodPut, url, Serialize(putNote))
	req.SetBasicAuth(token.Id, token.Token)
	resp := SendReq(req)
	ExpectStatusCode(resp, http.StatusOK)

	dbNote := getNote(noteId, m)
	ExpectEqual(*dbNote, *putNote)
}

func TestNotePutNoAuth(t *testing.T) {
	ts, _ := Server()
	defer ts.Close()
	url := buildNoteUrl("dogs", ts.URL)
	UnauthorizedTest(http.MethodPut, url, nil)
}
