package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/database"
	"notemeal-server/internal/handler"
	notemealModel "notemeal-server/internal/model"
)

func main() {
	db := database.DictDb()
	m := notemealModel.NewModel(db)
	mux := handler.ServeMux(m)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
