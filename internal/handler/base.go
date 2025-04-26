package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"notemeal-server/internal/database"
)

const objIdKey string = "id"

type baseHandler struct {
	Db          *database.Database
	Request     *http.Request
	RequestBody []byte
	Writer      http.ResponseWriter
	ObjId       string
	PrincipalId string
}

func (handler *baseHandler) endRequest(success bool) {
	if success {
		fmt.Printf("%s %s : success\n", handler.Request.Method, handler.Request.URL)
	} else {
		fmt.Printf("%s %s : failure\n", handler.Request.Method, handler.Request.URL)
	}
}

func (handler *baseHandler) getAuth() bool {
	var ok bool
	var tokenString string
	handler.PrincipalId, tokenString, ok = handler.Request.BasicAuth()

	if !ok {
		fmt.Println("Failed to find auth!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	user, err := (*handler.Db).GetUser(handler.PrincipalId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if user == nil {
		fmt.Println("Failed to find auth user in database!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	token, err := (*handler.Db).GetToken(tokenString)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if token == nil {
		fmt.Println("Failed to find auth token in database!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	if token.UserId != handler.PrincipalId {
		fmt.Println("Token does not belong to user!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func (handler *baseHandler) getBodyString() bool {
	var err error
	handler.RequestBody, err = io.ReadAll(handler.Request.Body)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not retrieve body string from Request!")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *baseHandler) getObjId() bool {
	handler.ObjId = handler.Request.PathValue(objIdKey)

	if handler.ObjId == "" {
		err := Error{"Could not get nodeId from path!"}
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *baseHandler) writeValueToResponse(value any) bool {
	body, err := json.Marshal(value)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	handler.Writer.Header().Add("Content-Type", "application/json")
	_, err = handler.Writer.Write(body)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}
