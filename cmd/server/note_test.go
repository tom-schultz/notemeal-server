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

func createNote(n *internal.Note) {
	err := database.Db.CreateNote(n)

	if err != nil {
		log.Fatal(err)
	}
}

func getNote(id string) *internal.Note {
	note, err := database.Db.GetNote(id)

	if err != nil {
		log.Fatal(err)
	}

	return note
}

func getNoteUrl(id string, baseUrl string) string {
	return fmt.Sprintf("%s/note/%s", baseUrl, id)
}

func TestNoteDeleteNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	url := getNoteUrl("dogs", ts.URL)
	test.UnauthorizedTest("DELETE", url, nil)
}

func TestNoteDelete(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	token := test.SetupAuth(user)

	note := getNote(noteId)

	if note == nil {
		log.Fatal("Test note does not exist!")
	}

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodDelete, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	note = getNote(noteId)

	if note != nil {
		log.Fatal("Note should be nil!")
	}
}

func TestNoteGetNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	url := getNoteUrl("dogs", ts.URL)
	test.UnauthorizedTest("GET", url, nil)
}

func TestNoteGet(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	token := test.SetupAuth(user)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodGet, url, nil)
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	note := getNote(noteId)
	test.ExpectBody(resp, test.Serialize(note))
}

func TestNotePutCreate(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "wuppers"
	postNote := &internal.Note{Id: noteId, Text: "woof woof", Title: "Puppers", UserId: user, LastModified: 0}
	token := test.SetupAuth(user)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodPut, url, test.Serialize(postNote))
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	dbNote := getNote(noteId)
	test.ExpectEqual(*dbNote, *postNote)
}

func TestNotePutUpdate(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	user := "tom"
	noteId := "dogs"
	putNote := &internal.Note{Id: noteId, Text: "woof woof", Title: "Puppers", UserId: user, LastModified: 0}
	token := test.SetupAuth(user)

	url := fmt.Sprintf("%s/note/%s", ts.URL, noteId)
	req := test.NewReq(http.MethodPut, url, test.Serialize(putNote))
	req.SetBasicAuth(token.Id, token.Token)
	resp := test.SendReq(req)
	test.ExpectStatusCode(resp, http.StatusOK)

	dbNote := getNote(noteId)
	test.ExpectEqual(*dbNote, *putNote)
}

func TestNotePutNoAuth(t *testing.T) {
	ts := test.Server()
	defer ts.Close()
	url := getNoteUrl("dogs", ts.URL)
	test.UnauthorizedTest(http.MethodPut, url, nil)
}
