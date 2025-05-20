package handler

import (
	"net/http"
	"notemeal-server/internal"
	"notemeal-server/internal/model"
)

func GetHealth(m model.Model, writer http.ResponseWriter, request *http.Request) {
	internal.LogRequestStart(request)
	handler := baseHandler{model: m, request: request, writer: writer}
	handler.writeValueToResponse("OK")
	handler.endRequest(true)
}
