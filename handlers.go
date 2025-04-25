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
