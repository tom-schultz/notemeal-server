package handler

import (
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
)

type codeHandler struct {
	baseHandler
	clientCode internal.ClientCode
}

func PutCode(m model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startCodeRequest(m, writer, request)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizeOwnsObj() &&
		handler.createTokenCode()

	handler.endRequest(result)
}

func (handler *codeHandler) createTokenCode() bool {
	code, err := handler.model.CreateOrUpdateCode(handler.objId)
	handler.clientCode = internal.ClientCode{UserId: handler.objId, Code: code}

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func startCodeRequest(m model.Model, writer http.ResponseWriter, request *http.Request) (*codeHandler, bool) {
	internal.LogRequestStart(request)

	handler := &codeHandler{
		baseHandler: baseHandler{
			model:   m,
			request: request,
			writer:  writer,
		},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}
