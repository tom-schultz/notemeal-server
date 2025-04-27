package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/database"
	"notemeal-server/internal/handler"
)

func main() {
	database.DictDb()
	mux := handler.ServeMux()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
