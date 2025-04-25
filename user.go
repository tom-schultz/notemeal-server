package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id          string
	ContactInfo string
}

type UserHandler struct {
	BaseHandler
	User *User
}

func (handler *UserHandler) deleteUser() bool {
	err := (*handler.Db).DeleteUser(handler.ObjId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *UserHandler) getUserFromDb() bool {
	var err error
	handler.User, err = (*handler.Db).GetUser(handler.ObjId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if handler.User == nil {
		handler.Writer.WriteHeader(http.StatusNotFound)
		return false
	}

	return true
}

func (handler *UserHandler) getUserFromBody() bool {
	handler.User = new(User)
	err := json.Unmarshal(handler.RequestBody, handler.User)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not deserialize note from body string!\n")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	handler.User.Id = handler.ObjId
	return true
}

func handleUserDELETE(writer http.ResponseWriter, request *http.Request) {
	handler := startUserRequest(writer, request, &db)

	result := handler.getObjId() &&
		handler.deleteUser()

	handler.endRequest(result)
}

func handleUserGET(writer http.ResponseWriter, request *http.Request) {
	handler := startUserRequest(writer, request, &db)

	result := handler.getObjId() &&
		handler.getUserFromDb() &&
		handler.writeValueToResponse(handler.User)

	handler.endRequest(result)
}

func handleUserPUT(writer http.ResponseWriter, request *http.Request) {
	handler := startUserRequest(writer, request, &db)

	result := handler.getObjId() &&
		handler.getBodyString() &&
		handler.getUserFromBody() &&
		handler.writeUserToDb()

	handler.endRequest(result)
}

func startUserRequest(writer http.ResponseWriter, request *http.Request, db *NotemealDb) *UserHandler {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)

	return &UserHandler{
		BaseHandler: BaseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}
}

func (handler *UserHandler) writeUserToDb() bool {
	if handler.ObjId != handler.User.Id {
		fmt.Println("Path objId does not match body objId!")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	err := (*handler.Db).SetUser(handler.User)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}
