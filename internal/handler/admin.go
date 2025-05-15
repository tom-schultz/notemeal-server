package handler

import (
	"fmt"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
)

type adminHandler struct {
	codeHandler
	code  string
	token string
}

func PutCodeAdmin(model model.Model, writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startAdminRequest(model, writer, request)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizeAdmin() &&
		handler.createTokenCode() &&
		handler.writeValueToResponse(handler.clientCode)

	handler.endRequest(result)
}

func (handler *adminHandler) authorizeAdmin() bool {
	authorized, err := handler.model.IsAdmin(handler.principalId)

	if err != nil {
		internal.LogRequestError(err, handler.request)
		handler.setStatus(http.StatusInternalServerError)
	}

	if !authorized {
		msg := fmt.Sprintf("%s is not an admin!\n", handler.principalId)
		internal.LogRequestMsg(msg, handler.request)
		handler.setStatus(http.StatusUnauthorized)
	}

	return authorized
}

func startAdminRequest(model model.Model, writer http.ResponseWriter, request *http.Request) (*adminHandler, bool) {
	internal.LogRequestStart(request)

	handler := &adminHandler{
		codeHandler: codeHandler{
			baseHandler: baseHandler{
				model:   model,
				request: request,
				writer:  writer,
			}},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}
