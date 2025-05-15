package handler

import (
	"net/http"
	"notemeal-server/internal/model"
)

type Handler struct {
	model       model.Model
	handlerFunc func(model.Model, http.ResponseWriter, *http.Request)
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.handlerFunc(h.model, writer, request)
}
