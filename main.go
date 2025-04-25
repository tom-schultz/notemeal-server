package main

import (
	"log"
	"net/http"
)

var db NotemealDb = new(NotemealDictDb)

func main() {
	if err := db.Initialize(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("DELETE /note/{id}", handleNoteDELETE)
	http.HandleFunc("GET /note/{id}", handleNoteGET)
	http.HandleFunc("PUT /note/{id}", handleNotePUT)
	http.HandleFunc("GET /notes", handleNotesGET)
	http.HandleFunc("DELETE /user/{id}", handleUserDELETE)
	http.HandleFunc("GET /user/{id}", handleUserGET)
	http.HandleFunc("PUT /user/{id}", handleUserPUT)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
