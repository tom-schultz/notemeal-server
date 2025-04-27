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
	handler, authenticated := startUserRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.deleteUser()

	handler.endRequest(result)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.getUserFromDb() &&
		handler.writeValueToResponse(handler.user)

	handler.endRequest(result)
}

func PutUser(writer http.ResponseWriter, request *http.Request) {
	handler, authenticated := startUserRequest(writer, request, &database.Db)

	if !authenticated {
		handler.endRequest(false)
		return
	}

	result := handler.getObjId() &&
		handler.authorizePrincipal() &&
		handler.getBodyString() &&
		handler.getUserFromBody() &&
		handler.writeUserToDb()

	handler.endRequest(result)
}

func (handler *userHandler) authorizePrincipal() bool {
	authorized := handler.PrincipalId == handler.ObjId

	if !authorized {
		fmt.Printf("%s is not authorized for operations on %s!\n", handler.PrincipalId, handler.ObjId)
		handler.Writer.WriteHeader(http.StatusUnauthorized)
	}

	return authorized
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

func startUserRequest(writer http.ResponseWriter, request *http.Request, db *database.Database) (*userHandler, bool) {
	fmt.Printf("%s %s : start\n", request.Method, request.URL)

	handler := &userHandler{
		baseHandler: baseHandler{
			Db:      db,
			Request: request,
			Writer:  writer,
		},
	}

	authenticated := handler.getAuth()
	return handler, authenticated
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
