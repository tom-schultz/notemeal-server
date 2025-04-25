package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Note struct {
	Id           string
	LastModified int
	Text         string
	Title        string
	UserId       string
}

type NoteHandler struct {
	BaseHandler
	Note     *Note
	NoteList map[string]int
}

func (handler *NoteHandler) authorizePrincipal() bool {
	isOwner, err := (*handler.Db).IsNoteOwner(handler.ObjId, handler.principalId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !isOwner {
		fmt.Println("Principal is not authorized!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func (handler *NoteHandler) deleteNote() bool {
	err := (*handler.Db).DeleteNote(handler.ObjId)

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

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not deserialize note from body string!\n")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	handler.Note.Id = handler.ObjId
	return true
}

func (handler *NoteHandler) getNoteFromDb() bool {
	var err error
	handler.Note, err = (*handler.Db).GetNote(handler.ObjId)

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

func handleNoteDELETE(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.deleteNote()

	handler.endRequest(result)
}

func handleNoteGET(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.getNoteFromDb() &&
		handler.writeValueToResponse(handler.Note)

	handler.endRequest(result)
}

func handleNotePUT(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.getBodyString() &&
		handler.getNoteFromBody() &&
		handler.updateNoteInDb()

	handler.endRequest(result)

}

func handleNotesGET(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.listNotesFromDb(handler.principalId) &&
		handler.writeValueToResponse(handler.NoteList)

	handler.endRequest(result)
}

func (handler *NoteHandler) listNotesFromDb(userId string) bool {
	var err error
	handler.NoteList, err = (*handler.Db).ListLastModified(userId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func startNoteRequest(writer http.ResponseWriter, request *http.Request, db *NotemealDb) (*NoteHandler, bool) {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)
	handler := &NoteHandler{
		BaseHandler: BaseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}
	authenticated := handler.getAuth()
	return handler, authenticated
}

func (handler *NoteHandler) updateNoteInDb() bool {
	if handler.ObjId != handler.Note.Id {
		fmt.Println("Path objId does not match body objId!")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	err := (*handler.Db).UpdateNote(handler.Note)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}
