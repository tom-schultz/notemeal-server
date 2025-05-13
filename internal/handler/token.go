package handler

import (
	"encoding/json"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type tokenHandler struct {
	baseHandler
	code        string
	clientToken *internal.ClientToken
}

func PostToken(writer http.ResponseWriter, request *http.Request) {
	handler := startTokenRequest(writer, request, &database.Db)

	result := handler.getObjId() &&
		handler.getBodyString() &&
		handler.getCodeFromBody() &&
		handler.createToken() &&
		handler.writeTokenToResponse()

	handler.endRequest(result)
}

func (handler *tokenHandler) createToken() bool {
	var err error
	handler.clientToken, err = (*handler.Db).CreateToken(handler.ObjId, handler.code)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if handler.clientToken == nil {
		internal.LogRequestMsg("Invalid code!", handler.Request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	return true
}

func (handler *tokenHandler) getCodeFromBody() bool {
	clientCode := internal.ClientCode{}
	err := json.Unmarshal(handler.RequestBody, &clientCode)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		internal.LogRequestMsg("Could not deserialize token code from body string!\n", handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	handler.code = clientCode.Code
	return true
}

func startTokenRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) *tokenHandler {
	internal.LogRequestStart(request)

	return &tokenHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}
}

func (handler *tokenHandler) writeTokenToResponse() bool {
	return handler.writeValueToResponse(*handler.clientToken)
}
