package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
)

const objIdKey string = "id"

type baseHandler struct {
	model       model.Model
	objId       string
	principalId string
	request     *http.Request
	requestBody []byte
	statusCode  int
	writer      http.ResponseWriter
}

func (handler *baseHandler) authorizeOwnsObj() bool {
	authorized := handler.principalId == handler.objId

	if !authorized {
		msg := fmt.Sprintf("%s is not authorized for operations on %s!\n", handler.principalId, handler.objId)
		internal.LogRequestMsg(msg, handler.request)
		handler.setStatus(http.StatusUnauthorized)
	}

	return authorized
}

func (handler *baseHandler) endRequest(success bool) {
	internal.LogRequestEnd(handler.request, handler.statusCode)
}

func (handler *baseHandler) authenticate() bool {
	var ok bool
	var tokenString string
	tokenId, tokenString, ok := handler.request.BasicAuth()

	if !ok {
		internal.LogRequestMsg("Failed to find auth!", handler.request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	principalToken, err := handler.model.GetToken(tokenId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if principalToken == nil {
		internal.LogRequestMsg("Failed to find token in database!", handler.request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	err = internal.CompareHashAndString(principalToken.Hash, tokenString)

	if err != nil {
		internal.LogRequestMsg("Invalid token!", handler.request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	handler.principalId = principalToken.UserId
	return true
}

func (handler *baseHandler) getBodyString() bool {
	var err error
	handler.requestBody, err = io.ReadAll(handler.request.Body)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		internal.LogRequestMsg("Could not retrieve body string from request!", handler.request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *baseHandler) getObjId() bool {
	handler.objId = handler.request.PathValue(objIdKey)

	if handler.objId == "" {
		err := Error{"Could not get nodeId from path!"}
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *baseHandler) setStatus(status int) {
	handler.statusCode = status
	handler.writer.WriteHeader(status)
}

func (handler *baseHandler) writeValueToResponse(value any) bool {
	body, err := json.Marshal(value)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	handler.writer.Header().Add("Content-Type", "application/json")
	_, err = handler.writer.Write(body)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}
