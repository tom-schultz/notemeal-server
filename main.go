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

	http.HandleFunc("DELETE /note/{noteId}", handleNoteDELETE)
	http.HandleFunc("GET /note/{noteId}", handleNoteGET)
	http.HandleFunc("PUT /note/{noteId}", handleNotePUT)
	http.HandleFunc("GET /notes", handleNotesGET)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
