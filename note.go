package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Note struct {
	Id           string
	Title        string
	Text         string
	LastModified int
}

const NoteIdKey = "noteId"

type NoteHandler struct {
	BaseHandler
	Note     *Note
	NoteList map[string]int
}

func (handler *NoteHandler) deleteNote() bool {
	err := (*handler.Db).DeleteNote(handler.PathValues[NoteIdKey])

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *NoteHandler) getNoteFromBody() bool {
	handler.Note = new(Note)
	err := json.Unmarshal(handler.RequestBody, handler.Note)
	handler.Note.Id = handler.PathValues[NoteIdKey]

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not deserialize note from body string!\n")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *NoteHandler) getNoteFromDb() bool {
	var err error
	handler.Note, err = (*handler.Db).GetNote(handler.PathValues[NoteIdKey])

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if handler.Note == nil {
		handler.Writer.WriteHeader(http.StatusNotFound)
		return false
	}

	return true
}

func (handler *NoteHandler) listNotesFromDb() bool {
	var err error
	handler.NoteList, err = (*handler.Db).ListLastModified()

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *NoteHandler) writeNoteToDb() bool {
	if handler.PathValues[NoteIdKey] != handler.Note.Id {
		fmt.Println("Path noteId does not match body noteId!")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	err := (*handler.Db).SetNote(handler.Note)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *NoteHandler) writeNoteToResponse() bool {
	body, err := json.Marshal(handler.Note)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	_, err = handler.Writer.Write(body)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *NoteHandler) writeNotesToResponse() bool {
	respBody, err := json.Marshal(handler.NoteList)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	_, err = handler.Writer.Write(respBody)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func handleNoteDELETE(writer http.ResponseWriter, request *http.Request) {
	handler := startRequest(writer, request, &db)

	result := handler.getPathValue(NoteIdKey) &&
		handler.deleteNote()

	handler.endRequest(result)
}

func handleNoteGET(writer http.ResponseWriter, request *http.Request) {
	handler := startRequest(writer, request, &db)

	result := handler.getPathValue(NoteIdKey) &&
		handler.getNoteFromDb() &&
		handler.writeNoteToResponse()

	handler.endRequest(result)
}

func handleNotePUT(writer http.ResponseWriter, request *http.Request) {
	handler := startRequest(writer, request, &db)

	result := handler.getPathValue(NoteIdKey) &&
		handler.getBodyString() &&
		handler.getNoteFromBody() &&
		handler.writeNoteToDb()

	handler.endRequest(result)
}

func handleNotesGET(writer http.ResponseWriter, request *http.Request) {
	handler := startRequest(writer, request, &db)

	result := handler.listNotesFromDb() &&
		handler.writeNotesToResponse()

	handler.endRequest(result)
}
