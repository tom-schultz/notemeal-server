package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

const objIdKey string = "id"

type baseHandler struct {
	Db          *database.Database
	ObjId       string
	PrincipalId string
	Request     *http.Request
	RequestBody []byte
	StatusCode  int
	Writer      http.ResponseWriter
}

func (handler *baseHandler) authorizePrincipal() bool {
	authorized := handler.PrincipalId == handler.ObjId

	if !authorized {
		msg := fmt.Sprintf("%s is not authorized for operations on %s!\n", handler.PrincipalId, handler.ObjId)
		internal.LogRequestMsg(msg, handler.Request)
		handler.setStatus(http.StatusUnauthorized)
	}

	return authorized
}

func (handler *baseHandler) endRequest(success bool) {
	internal.LogRequestEnd(handler.Request, handler.StatusCode)
}

func (handler *baseHandler) getAuth() bool {
	var ok bool
	var tokenString string
	handler.PrincipalId, tokenString, ok = handler.Request.BasicAuth()

	if !ok {
		internal.LogRequestMsg("Failed to find auth!", handler.Request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	user, err := (*handler.Db).GetUser(handler.PrincipalId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if user == nil {
		internal.LogRequestMsg("Failed to find auth user in database!", handler.Request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	token, err := (*handler.Db).GetToken(tokenString)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if token == nil {
		internal.LogRequestMsg("Failed to find auth token in database!", handler.Request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	if token.UserId != handler.PrincipalId {
		internal.LogRequestMsg("Token does not belong to user!", handler.Request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	return true
}

func (handler *baseHandler) getBodyString() bool {
	var err error
	handler.RequestBody, err = io.ReadAll(handler.Request.Body)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		internal.LogRequestMsg("Could not retrieve body string from Request!", handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *baseHandler) getObjId() bool {
	handler.ObjId = handler.Request.PathValue(objIdKey)

	if handler.ObjId == "" {
		err := Error{"Could not get nodeId from path!"}
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *baseHandler) setStatus(status int) {
	handler.StatusCode = status
	handler.Writer.WriteHeader(status)
}

func (handler *baseHandler) writeValueToResponse(value any) bool {
	body, err := json.Marshal(value)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	handler.Writer.Header().Add("Content-Type", "application/json")
	_, err = handler.Writer.Write(body)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}
