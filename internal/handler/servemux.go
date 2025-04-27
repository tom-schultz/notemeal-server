package handler

import (
	"net/http"
)

func ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /note/{id}", DeleteNote)
	mux.HandleFunc("GET /note/{id}", GetNote)
	mux.HandleFunc("PUT /note/{id}", PutNote)
	mux.HandleFunc("GET /notes", GetNotes)
	mux.HandleFunc("DELETE /user/{id}", DeleteUser)
	mux.HandleFunc("GET /user/{id}", GetUser)
	mux.HandleFunc("PUT /user/{id}", PutUser)
	mux.HandleFunc("PUT /user/{id}/code", PutCode)
	mux.HandleFunc("POST /user/{id}/token", PostToken)

	return mux
}
