package handler

import (
	"net/http"
	notemealModel "notemeal-server/internal/model"
)

func ServeMux(m notemealModel.Model) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("PUT /admin/user/{id}/code", Handler{m, PutCodeAdmin})
	mux.Handle("DELETE /note/{id}", Handler{m, DeleteNote})
	mux.Handle("GET /note/{id}", Handler{m, GetNote})
	mux.Handle("PUT /note/{id}", Handler{m, PutNote})
	mux.Handle("GET /notes", Handler{m, GetNotes})
	mux.Handle("DELETE /user/{id}", Handler{m, DeleteUser})
	mux.Handle("GET /user/{id}", Handler{m, GetUser})
	mux.Handle("PUT /user/{id}", Handler{m, PutUser})
	mux.Handle("PUT /user/{id}/code", Handler{m, PutCode})
	mux.Handle("POST /user/{id}/token", Handler{m, PostToken})
	mux.Handle("GET /health", Handler{m, GetHealth})

	return mux
}
