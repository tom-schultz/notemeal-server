package handler

import (
	"encoding/json"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
)

type userHandler struct {
	baseHandler
	user *internal.User
}

func DeleteUser(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(m, writer, request)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizeOwnsObj() &&
		handler.deleteUser()

	handler.endRequest(result)
}

func GetUser(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(m, writer, request)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizeOwnsObj() &&
		handler.getUserFromDb() &&
		handler.writeValueToResponse(handler.user)

	handler.endRequest(result)
}

func PutUser(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(m, writer, request)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizeOwnsObj() &&
		handler.getBodyString() &&
		handler.getUserFromBody() &&
		handler.writeUserToDb()

	handler.endRequest(result)
}

func (handler *userHandler) deleteUser() bool {
	err := handler.model.DeleteUser(handler.objId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *userHandler) getUserFromDb() bool {
	var err error
	handler.user, err = handler.model.GetUser(handler.objId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if handler.user == nil {
		handler.setStatus(http.StatusNotFound)
		return false
	}

	return true
}

func (handler *userHandler) getUserFromBody() bool {
	handler.user = new(internal.User)
	err := json.Unmarshal(handler.requestBody, handler.user)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		internal.LogRequestMsg("Could not deserialize user from body string!\n", handler.request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	handler.user.Id = handler.objId
	return true
}

func startUserRequest(m model.Model, writer http.ResponseWriter, request *http.Request) (*userHandler, bool) {
	internal.LogRequestStart(request)

	handler := &userHandler{
		baseHandler: baseHandler{
			model:   m,
			request: request,
			writer:  writer,
		},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}

func (handler *userHandler) writeUserToDb() bool {
	if handler.objId != handler.user.Id {
		internal.LogRequestMsg("Path objId does not match body objId!", handler.request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	err := handler.model.SetUser(handler.user)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}
