package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type noteHandler struct {
	baseHandler
	note     *internal.Note
	noteList map[string]int
}

func DeleteNote(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.deleteNote()

	handler.endRequest(result)
}

func GetNote(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.getNoteFromDb() &&
		handler.writeValueToResponse(handler.note)

	handler.endRequest(result)
}

func PutNote(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &database.Db)

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

func GetNotes(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.listNotesFromDb(handler.PrincipalId) &&
		handler.writeValueToResponse(handler.noteList)

	handler.endRequest(result)
}

func (handler *noteHandler) authorizePrincipal() bool {
	isOwner, err := (*handler.Db).IsNoteOwner(handler.ObjId, handler.PrincipalId)

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

func (handler *noteHandler) deleteNote() bool {
	err := (*handler.Db).DeleteNote(handler.ObjId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *noteHandler) getNoteFromBody() bool {
	handler.note = new(internal.Note)
	err := json.Unmarshal(handler.RequestBody, handler.note)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not deserialize note from body string!\n")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	handler.note.Id = handler.ObjId
	return true
}

func (handler *noteHandler) getNoteFromDb() bool {
	var err error
	handler.note, err = (*handler.Db).GetNote(handler.ObjId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if handler.note == nil {
		handler.Writer.WriteHeader(http.StatusNotFound)
		return false
	}

	return true
}

func (handler *noteHandler) listNotesFromDb(userId string) bool {
	var err error
	handler.noteList, err = (*handler.Db).ListLastModified(userId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func startNoteRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) (*noteHandler, bool) {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)
	handler := &noteHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}
	authenticated := handler.getAuth()
	return handler, authenticated
}

func (handler *noteHandler) updateNoteInDb() bool {
	if handler.ObjId != handler.note.Id {
		fmt.Println("Path objId does not match body objId!")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	err := (*handler.Db).UpdateNote(handler.note)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}
