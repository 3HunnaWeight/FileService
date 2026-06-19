package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

//type Handler struct {
//	fileHandler *FileHandler
//}

func NewRouter(fileHandler *FileHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/files", fileHandler.Upload)
	r.Get("/files/{id}", fileHandler.Download)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
