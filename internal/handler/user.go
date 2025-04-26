package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/database"
)

type userHandler struct {
	baseHandler
	user *internal.User
}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	handler := startUserRequest(writer, request, &database.Db)

	result := handler.getObjId() &&
		handler.deleteUser()

	handler.endRequest(result)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	handler := startUserRequest(writer, request, &database.Db)

	result := handler.getObjId() &&
		handler.getUserFromDb() &&
		handler.writeValueToResponse(handler.user)

	handler.endRequest(result)
}

func PutUser(writer http.ResponseWriter, request *http.Request) {
	handler := startUserRequest(writer, request, &database.Db)

	result := handler.getObjId() &&
		handler.getBodyString() &&
		handler.getUserFromBody() &&
		handler.writeUserToDb()

	handler.endRequest(result)
}

func (handler *userHandler) deleteUser() bool {
	err := (*handler.Db).DeleteUser(handler.ObjId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func (handler *userHandler) getUserFromDb() bool {
	var err error
	handler.user, err = (*handler.Db).GetUser(handler.ObjId)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if handler.user == nil {
		handler.Writer.WriteHeader(http.StatusNotFound)
		return false
	}

	return true
}

func (handler *userHandler) getUserFromBody() bool {
	handler.user = new(internal.User)
	err := json.Unmarshal(handler.RequestBody, handler.user)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not deserialize user from body string!\n")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	handler.user.Id = handler.ObjId
	return true
}

func startUserRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) *userHandler {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)

	return &userHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}
}

func (handler *userHandler) writeUserToDb() bool {
	if handler.ObjId != handler.user.Id {
		fmt.Println("Path objId does not match body objId!")
		handler.Writer.WriteHeader(http.StatusBadRequest)
		return false
	}

	err := (*handler.Db).SetUser(handler.user)

	if err != nil {
		fmt.Println(err)
		handler.Writer.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}
