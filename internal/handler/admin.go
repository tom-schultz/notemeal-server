package handler

import (
	"fmt"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type adminHandler struct {
	codeHandler
	code  string
	token string
}

func PutCodeAdmin(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startAdminRequest(writer, request, &database.Db)

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
	authorized, err := (*handler.Db).IsAdmin(handler.PrincipalId)

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
	}

	if !authorized {
		msg := fmt.Sprintf("%s is not an admin!\n", handler.PrincipalId)
		internal.LogRequestMsg(msg, handler.Request)
		handler.setStatus(http.StatusUnauthorized)
	}

	return authorized
}

func startAdminRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) (*adminHandler, bool) {
	internal.LogRequestStart(request)

	handler := &adminHandler{
		codeHandler: codeHandler{
			baseHandler: baseHandler{
				Db:      db,
				Request: request,
				Writer:  writer,
			}},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}
