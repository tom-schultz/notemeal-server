package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
)

type Note struct {
	Id           string
	Title        string
	Text         string
	LastModified int
}

var notes = map[string]*Note{
	"doggos":  {"doggos", "Doggos", "doggos are sweet", 0},
	"cattos":  {"cattos", "Cattos", "meowow", 0},
	"rabbits": {"rabbits", "Rabbits", "hoppity hop, mothafucka", 0},
}

func main() {
	http.HandleFunc("DELETE /note/{noteId}", handleNoteDELETE)
	http.HandleFunc("GET /note/{noteId}", handleNoteGET)
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
	note, ok := notes[noteId]

	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprintf(writer, "There is no note with ID %s", noteId)

		if err != nil {
			fmt.Println("Got an error!")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	body, err := json.Marshal(note)

	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(body)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func handleNotePUT(writer http.ResponseWriter, request *http.Request) {
	noteId := request.PathValue("noteId")
	respBody := "OK"

	if _, ok := notes[noteId]; !ok {
		fmt.Printf("Note not found, creating new note %s!\n", noteId)
		notes[noteId] = new(Note)
	}

	reqBody, err := io.ReadAll(request.Body)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not retrieve body from request for note %s!\n", noteId)
		writer.WriteHeader(http.StatusBadRequest)
		respBody = "MALFORMED REQUEST"
	} else {
		err = json.Unmarshal(reqBody, notes[noteId])

		if err != nil {
			fmt.Println(err)
			fmt.Printf("Could not parse body json in request for note %s!\n", noteId)
			writer.WriteHeader(http.StatusBadRequest)
			respBody = "MALFORMED REQUEST"
		}
	}

	_, err = fmt.Fprintf(writer, respBody)

	if err != nil {
		fmt.Println(err)
	}
}

func handleNotesGET(writer http.ResponseWriter, request *http.Request) {
	respData := make(map[string]int)

	for key := range maps.Keys(notes) {
		respData[key] = notes[key].LastModified
	}

	respBody, err := json.Marshal(respData)

	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}

	_, err = writer.Write(respBody)

	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
