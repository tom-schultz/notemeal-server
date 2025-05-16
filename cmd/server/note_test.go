package main

import (
	"fmt"
	"log"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
	"notemeal-server/internal/test"
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
	ts, _ := test.Server()
	defer ts.Close()
	url := buildNoteUrl("dogs", ts.URL)
	test.UnauthorizedTest("DELETE", url, nil)
}

func TestNoteDelete(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	token := test.SetupAuth(user, m)

	note := getNote(noteId, m)

	if note == nil {
		log.Fatal("Test note does not exist!")
	}

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	note = getNote(noteId, m)

	if note != nil {
		log.Fatal("Note should be nil!")
	}
}

func TestNoteGetNoAuth(t *testing.T) {
	ts, _ := test.Server()
	defer ts.Close()
	url := buildNoteUrl("dogs", ts.URL)
	test.UnauthorizedTest("GET", url, nil)
}

func TestNoteGet(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	token := test.SetupAuth(user, m)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	note := getNote(noteId, m)
	test.ExpectBody(resp, test.Serialize(note))
}

func TestNotePutCreate(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "wuppers"
	postNote := &internal.Note{Id: noteId, Text: "woof woof", Title: "Puppers", UserId: user, LastModified: 0}
	token := test.SetupAuth(user, m)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodPut, url, test.Serialize(postNote))
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	dbNote := getNote(noteId, m)
	test.ExpectEqual(*dbNote, *postNote)
}

func TestNotePutUpdate(t *testing.T) {
	ts, m := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	putNote := &internal.Note{Id: noteId, Text: "woof woof", Title: "Puppers", UserId: user, LastModified: 0}
	token := test.SetupAuth(user, m)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodPut, url, test.Serialize(putNote))
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	dbNote := getNote(noteId, m)
	test.ExpectEqual(*dbNote, *putNote)
}

func TestNotePutNoAuth(t *testing.T) {
	ts, _ := test.Server()
	defer ts.Close()
	url := buildNoteUrl("dogs", ts.URL)
	test.UnauthorizedTest(http.MethodPut, url, nil)
}
