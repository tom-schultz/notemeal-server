package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"slices"
	"strings"
)

var notes = map[string]string{
	"doggos":  "doggos are sweet",
	"cattos":  "meowow",
	"rabbits": "hoppity hop, mothafucka",
}

func main() {
	http.HandleFunc("DELETE /note/{noteId}", handleNoteDELETE)
	http.HandleFunc("GET /note/{noteId}", handleNoteGET)
	http.HandleFunc("POST /note", handleNotePOST)
	http.HandleFunc("PUT /note/{noteId}", handleNotePUT)
	http.HandleFunc("GET /notes", handleNotesGET)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleNoteDELETE(writer http.ResponseWriter, request *http.Request) {
	noteId := request.PathValue("noteId")
	respBody := "OK"

	if _, ok := notes[noteId]; !ok {
		return
	}

	delete(notes, noteId)
	_, err := fmt.Fprintf(writer, respBody)

	if err != nil {
		fmt.Println(err)
	}
}

func handleNoteGET(writer http.ResponseWriter, request *http.Request) {
	noteId := request.PathValue("noteId")
	body, ok := notes[noteId]

	if !ok {
		body = fmt.Sprintf("There is no note with ID %s", noteId)
	}

	_, err := fmt.Fprintf(writer, body)

	if err != nil {
		fmt.Println(err)
	}
}

func handleNotePOST(writer http.ResponseWriter, request *http.Request) {
	reqBody, err := io.ReadAll(request.Body)

	if err != nil {
		fmt.Println(err)
		reqBody = make([]byte, 0)
	}

	noteId := pseudo_uuid()
	notes[noteId] = string(reqBody)

	_, err = fmt.Fprintf(writer, noteId)

	if err != nil {
		fmt.Println(err)
	}
}

func handleNotePUT(writer http.ResponseWriter, request *http.Request) {
	noteId := request.PathValue("noteId")
	respBody := "OK"

	if _, ok := notes[noteId]; !ok {
		fmt.Printf("Could not find note %s!\n", noteId)

		writer.WriteHeader(http.StatusNotFound)
		respBody = "NOTE NOT FOUND"
	} else {
		reqBody, err := io.ReadAll(request.Body)

		if err != nil {
			fmt.Println(err)
			fmt.Printf("Could note retrieve body from request for note %s!\n", noteId)
			writer.WriteHeader(http.StatusBadRequest)
			respBody = "MALFORMED REQUESTE"
		} else {
			notes[noteId] = string(reqBody)
		}
	}

	_, err := fmt.Fprintf(writer, respBody)

	if err != nil {
		fmt.Println(err)
	}
}

func handleNotesGET(writer http.ResponseWriter, request *http.Request) {
	respBody := strings.Join(slices.Collect(maps.Keys(notes)), "\n")

	_, err := fmt.Fprintf(writer, respBody)

	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// Note - NOT RFC4122 compliant
func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
