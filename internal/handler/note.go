package handler

import (
	"encoding/json"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type noteHandler struct {
	baseHandler
	note         *internal.Note
	noteList     map[string]int
	existingNote bool
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
	oldNote, err := (*handler.Db).GetNote(handler.ObjId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	handler.existingNote = oldNote != nil

	if handler.existingNote {
		isOwner, err := (*handler.Db).IsNoteOwner(handler.ObjId, handler.PrincipalId)

		if err != nil {
			internal.LogRequestError(err, handler.Request)
			handler.setStatus(http.StatusInternalServerError)
			return false
		}

		if !isOwner {
			internal.LogRequestMsg("Principal is not authorized!", handler.Request)
			handler.setStatus(http.StatusUnauthorized)
			return false
		}
	}

	return true
}

func (handler *noteHandler) deleteNote() bool {
	err := (*handler.Db).DeleteNote(handler.ObjId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *noteHandler) getNoteFromBody() bool {
	handler.note = new(internal.Note)
	err := json.Unmarshal(handler.RequestBody, handler.note)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		internal.LogRequestMsg("Could not deserialize note from body string!\n", handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	handler.note.Id = handler.ObjId
	return true
}

func (handler *noteHandler) getNoteFromDb() bool {
	var err error
	handler.note, err = (*handler.Db).GetNote(handler.ObjId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if handler.note == nil {
		handler.setStatus(http.StatusNotFound)
		return false
	}

	return true
}

func (handler *noteHandler) listNotesFromDb(userId string) bool {
	var err error
	handler.noteList, err = (*handler.Db).ListLastModified(userId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func startNoteRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) (*noteHandler, bool) {
	internal.LogRequestStart(request)

	handler := &noteHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}

func (handler *noteHandler) updateNoteInDb() bool {
	if handler.ObjId != handler.note.Id {
		internal.LogRequestMsg("Path objId does not match body objId!", handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	var err error
	handler.note.UserId = handler.PrincipalId

	if handler.existingNote {
		err = (*handler.Db).UpdateNote(handler.note)
	} else {
		err = (*handler.Db).CreateNote(handler.note)
	}

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}
