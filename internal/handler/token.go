package handler

import (
	"encoding/json"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type tokenHandler struct {
	baseHandler
	code  string
	token string
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
	handler.token, err = (*handler.Db).CreateToken(handler.ObjId, handler.code)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if handler.token == "" {
		internal.LogRequestMsg("Invalid code!", handler.Request)
		handler.setStatus(http.StatusUnauthorized)
	}

	return true
}

func (handler *tokenHandler) getCodeFromBody() bool {
	data := map[string]string{}
	err := json.Unmarshal(handler.RequestBody, &data)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		internal.LogRequestMsg("Could not deserialize token code from body string!\n", handler.Request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	handler.code = data[internal.CodeJsonKey]
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
	data := map[string]string{internal.TokenJsonKey: handler.token}
	return handler.writeValueToResponse(data)
}
