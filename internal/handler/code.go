package handler

import (
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type codeHandler struct {
	baseHandler
	clientCode internal.ClientCode
}

func PutCode(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startCodeRequest(writer, request, &database.Db)

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
	code, err := (*handler.Db).CreateOrUpdateCode(handler.ObjId)
	handler.clientCode = internal.ClientCode{UserId: handler.ObjId, Code: code}

	if err != nil {
		internal.LogRequestError(err, handler.Request)
		handler.setStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func startCodeRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) (*codeHandler, bool) {
	internal.LogRequestStart(request)

	handler := &codeHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}

	authenticated := handler.authenticate()
	return handler, authenticated
}
