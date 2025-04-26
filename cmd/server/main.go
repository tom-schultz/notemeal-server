package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/database"
	"notemeal-server/internal/handler"
)

func main() {
	database.DictDb()

	http.HandleFunc("DELETE /note/{id}", handler.DeleteNote)
	http.HandleFunc("GET /note/{id}", handler.GetNote)
	http.HandleFunc("PUT /note/{id}", handler.PutNote)
	http.HandleFunc("GET /notes", handler.GetNotes)
	http.HandleFunc("DELETE /user/{id}", handler.DeleteUser)
	http.HandleFunc("GET /user/{id}", handler.GetUser)
	http.HandleFunc("PUT /user/{id}", handler.PutUser)
	http.HandleFunc("PUT /user/{id}/code", handler.PutCode)
	http.HandleFunc("POST /user/{id}/token", handler.PostToken)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
