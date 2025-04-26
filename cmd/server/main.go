package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/database"
	"notemeal-server/internal/handler"
)

func main() {
	database.DictDb()
	mux := ServeMux()
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /note/{id}", handler.DeleteNote)
	mux.HandleFunc("GET /note/{id}", handler.GetNote)
	mux.HandleFunc("PUT /note/{id}", handler.PutNote)
	mux.HandleFunc("GET /notes", handler.GetNotes)
	mux.HandleFunc("DELETE /user/{id}", handler.DeleteUser)
	mux.HandleFunc("GET /user/{id}", handler.GetUser)
	mux.HandleFunc("PUT /user/{id}", handler.PutUser)
	mux.HandleFunc("PUT /user/{id}/code", handler.PutCode)
	mux.HandleFunc("POST /user/{id}/token", handler.PostToken)

	return mux
}
