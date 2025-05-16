package main

import (
	"log"
	"net/http"
	"notemeal-server/internal/data"
	"notemeal-server/internal/handler"
	notemealModel "notemeal-server/internal/model"
)

func main() {
	ds := data.DictDb()
	m := notemealModel.NewModel(ds)
	mux := handler.ServeMux(m)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
