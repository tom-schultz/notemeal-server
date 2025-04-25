package main

import (
	"fmt"
	"io"
	"net/http"
)

type HandlerError struct {
	msg string
}

func (e HandlerError) Error() string {
	return e.msg
}

type BaseHandler struct {
	Db          *NotemealDb
	PathValues  map[string]string
	Request     *http.Request
	RequestBody []byte
	Writer      http.ResponseWriter
}

func (handler *BaseHandler) endRequest(success bool) {
	if success {
		fmt.Printf("%s %s : success\n", handler.Request.Method, handler.Request.URL)
	} else {
		fmt.Printf("%s %s : failure\n", handler.Request.Method, handler.Request.URL)
	}
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

func (handler *BaseHandler) getPathValue(key string) bool {
	handler.PathValues[key] = handler.Request.PathValue(key)

	if handler.PathValues[key] == "" {
		err := HandlerError{"Could not get nodeId from path!"}
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func startRequest(writer http.ResponseWriter, request *http.Request, db *NotemealDb) *NoteHandler {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)

	return &NoteHandler{
		BaseHandler: BaseHandler{
			Db:         db,
			PathValues: make(map[string]string),
			Request:    request,
			Writer:     writer,
		},
	}
}
