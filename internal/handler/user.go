package handler

import (
	"encoding/json"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type userHandler struct {
	baseHandler
	user *internal.User
}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizeOwnsObj() &&
		handler.deleteUser()

	handler.endRequest(result)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(writer, request, &database.Db)

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

func PutUser(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(writer, request, &database.Db)

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
	err := (*handler.Db).DeleteUser(handler.ObjId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *userHandler) getUserFromDb() bool {
	var err error
	handler.user, err = (*handler.Db).GetUser(handler.ObjId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
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
	err := json.Unmarshal(handler.RequestBody, handler.user)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		internal.LogRequestMsg("Could not deserialize user from body string!\n", handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	handler.user.Id = handler.ObjId
	return true
}

func startUserRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) (*userHandler, bool) {
	internal.LogRequestStart(request)

	handler := &userHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}

func (handler *userHandler) writeUserToDb() bool {
	if handler.ObjId != handler.user.Id {
		internal.LogRequestMsg("Path objId does not match body objId!", handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	err := (*handler.Db).SetUser(handler.user)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}
