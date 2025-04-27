package handler

import (
	"encoding/json"
	"fmt"
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
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if handler.token == "" {
		fmt.Println("Invalid code!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
	}

	return true
}

func (handler *tokenHandler) getCodeFromBody() bool {
	data := map[string]string{}
	err := json.Unmarshal(handler.RequestBody, &data)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not deserialize token code from body string!\n")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	handler.code = data[internal.CodeJsonKey]
	return true
}

func startTokenRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) *tokenHandler {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)

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
