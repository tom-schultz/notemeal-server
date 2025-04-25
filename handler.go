package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ObjIdKey string = "id"

type HandlerError struct {
	msg string
}

func (e HandlerError) Error() string {
	return e.msg
}

type BaseHandler struct {
	Db          *NotemealDb
	Request     *http.Request
	RequestBody []byte
	Writer      http.ResponseWriter
	ObjId       string
	principalId string
}

func (handler *BaseHandler) endRequest(success bool) {
	if success {
		fmt.Printf("%s %s : success\n", handler.Request.Method, handler.Request.URL)
	} else {
		fmt.Printf("%s %s : failure\n", handler.Request.Method, handler.Request.URL)
	}
}

func (handler *BaseHandler) getAuth() bool {
	var ok bool
	var tokenString string
	handler.principalId, tokenString, ok = handler.Request.BasicAuth()

	if !ok {
		fmt.Println("Failed to find auth!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	user, err := (*handler.Db).GetUser(handler.principalId)

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

	if token.UserId != handler.principalId {
		fmt.Println("Token does not belong to user!")
		handler.Writer.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func (handler *BaseHandler) getBodyString() bool {
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

func (handler *BaseHandler) getObjId() bool {
	handler.ObjId = handler.Request.PathValue(ObjIdKey)

	if handler.ObjId == "" {
		err := HandlerError{"Could not get nodeId from path!"}
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (handler *BaseHandler) writeValueToResponse(value any) bool {
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
