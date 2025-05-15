package handler

import (
	"encoding/json"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
)

type tokenHandler struct {
	baseHandler
	code        string
	clientToken *internal.ClientToken
}

func PostToken(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler := startTokenRequest(m, writer, request)

	result := handler.getObjId() &&
		handler.getBodyString() &&
		handler.getCodeFromBody() &&
		handler.createToken() &&
		handler.writeTokenToResponse()

	handler.endRequest(result)
}

func (handler *tokenHandler) createToken() bool {
	var err error
	handler.clientToken, err = handler.model.CreateToken(handler.objId, handler.code)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	if handler.clientToken == nil {
		internal.LogRequestMsg("Invalid code!", handler.request)
		handler.setStatus(http.StatusUnauthorized)
		return false
	}

	return true
}

func (handler *tokenHandler) getCodeFromBody() bool {
	clientCode := internal.ClientCode{}
	err := json.Unmarshal(handler.requestBody, &clientCode)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		internal.LogRequestMsg("Could not deserialize token code from body string!\n", handler.request)
		handler.setStatus(http.StatusBadRequest)
		return false
	}

	handler.code = clientCode.Code
	return true
}

func startTokenRequest(m model.Model, writer http.ResponseWriter, request *http.Request) *tokenHandler {
	internal.LogRequestStart(request)

	return &tokenHandler{
		baseHandler: baseHandler{
			model:   m,
			request: request,
			writer:  writer,
		},
	}
}

func (handler *tokenHandler) writeTokenToResponse() bool {
	return handler.writeValueToResponse(*handler.clientToken)
}
