package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/3HunnaWeight/file-service/internal/domain"
	"github.com/go-chi/chi/v5"
)

type FileService interface {
	Upload(
		ctx context.Context,
		fileName string,
		data []byte,
		mime string,
	) (string, error)
	GetByPublicID(ctx context.Context, publicID string) (*domain.File, error)
	Download(
		ctx context.Context,
		publicID string,
	) (io.ReadCloser, *domain.File, error)
}

type FileHandler struct {
	service FileService
}

func NewFileHandler(service FileService) *FileHandler {
	return &FileHandler{service: service}
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publicID, err := h.service.Upload(
		r.Context(),
		header.Filename,
		data,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":  publicID,
		"url": "/files/" + publicID,
	})
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	publicID := chi.URLParam(r, "id")

	reader, file, err := h.service.Download(r.Context(), publicID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+file.OriginalName+`"`)

	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "stream failed", http.StatusInternalServerError)
	}
}
