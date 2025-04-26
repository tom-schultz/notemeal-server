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
	handler := startCodeRequest(writer, request, &database.Db)

	result := handler.getObjId() &&
		handler.createTokenCode()

	handler.endRequest(result)
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

func startCodeRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) *codeHandler {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)

	return &codeHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}
}
