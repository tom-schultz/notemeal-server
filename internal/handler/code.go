package handler

import (
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type codeHandler struct {
	baseHandler
	codeData map[string]string
}

func PutCode(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startCodeRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.createTokenCode()

	handler.endRequest(result)
}

func (handler *codeHandler) createTokenCode() bool {
	code, err := (*handler.Db).CreateOrUpdateCode(handler.ObjId)
	handler.codeData = map[string]string{internal.CodeJsonKey: code}

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

	authenticated := handler.getAuth()
	return handler, authenticated
}
