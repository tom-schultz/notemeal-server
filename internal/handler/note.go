package handler

import (
	"encoding/json"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
)

type noteHandler struct {
	baseHandler
	note         *internal.Note
	noteList     map[string]int
	existingNote bool
}

func DeleteNote(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(m, writer, request)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.deleteNote()

	handler.endRequest(result)
}

func GetNote(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(m, writer, request)

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

func PutNote(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(m, writer, request)

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

func GetNotes(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startNoteRequest(m, writer, request)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.listNotesFromDb(handler.principalId) &&
		handler.writeValueToResponse(handler.noteList)

	handler.endRequest(result)
}

func (handler *noteHandler) authorizePrincipal() bool {
	oldNote, err := handler.model.GetNote(handler.objId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	handler.existingNote = oldNote != nil

	if handler.existingNote {
		isOwner, err := handler.model.IsNoteOwner(handler.objId, handler.principalId)

		if err != nil {
			internal.LogRequestError(err, handler.request)
			handler.setStatus(http.StatusInternalServerError)
			return false
		}

		if !isOwner {
			internal.LogRequestMsg("Principal is not authorized!", handler.request)
			handler.setStatus(http.StatusUnauthorized)
			return false
		}
	}

	return true
}

func (handler *noteHandler) deleteNote() bool {
	err := handler.model.DeleteNote(handler.objId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *noteHandler) getNoteFromBody() bool {
	handler.note = new(internal.Note)
	err := json.Unmarshal(handler.requestBody, handler.note)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		internal.LogRequestMsg("Could not deserialize note from body string!\n", handler.request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	handler.note.Id = handler.objId
	return true
}

func (handler *noteHandler) getNoteFromDb() bool {
	var err error
	handler.note, err = handler.model.GetNote(handler.objId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
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
	handler.noteList, err = handler.model.ListLastModified(userId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func startNoteRequest(model model.Model, writer http.ResponseWriter, request *http.Request) (*noteHandler, bool) {
	internal.LogRequestStart(request)

	handler := &noteHandler{
		baseHandler: baseHandler{
			model:   model,
			request: request,
			writer:  writer,
		},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}

func (handler *noteHandler) updateNoteInDb() bool {
	if handler.objId != handler.note.Id {
		internal.LogRequestMsg("Path objId does not match body objId!", handler.request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	var err error
	handler.note.UserId = handler.principalId

	if handler.existingNote {
		err = handler.model.UpdateNote(handler.note)
	} else {
		err = handler.model.CreateNote(handler.note)
	}

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}
