package handler

import (
	"fmt"
	"net/http"
	"notemeal-server/internal/database"
)

type codeHandler struct {
	baseHandler
	code string
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

func (handler *codeHandler) authorizePrincipal() bool {
	authorized := handler.PrincipalId == handler.ObjId

	if !authorized {
		fmt.Printf("%s is not authorized for operations on %s!\n", handler.PrincipalId, handler.ObjId)
		handler.Writer.WriteHeader(http.StatusUnauthorized)
	}

	return authorized
}

func (handler *codeHandler) createTokenCode() bool {
	_, err := (*handler.Db).CreateOrUpdateCode(handler.ObjId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func startCodeRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) (*codeHandler, bool) {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)

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
