package http

import (
	"context"
	"encoding/json"
	"github.com/3HunnaWeight/file-service/internal/domain"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
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

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	publicID, err := h.service.Upload(
		r.Context(),
		header.Filename,
		data,
		header.Header.Get("Content-Type"),
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"id":  publicID,
		"url": "/f/" + publicID,
	})
}
func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	publicID := chi.URLParam(r, "id")

	file, err := h.service.GetByPublicID(r.Context(), publicID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	reader, file, err := h.service.Download(
		r.Context(),
		publicID,
	)
	if err != nil {
		http.Error(w, "download failed", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set(
		"Content-Disposition",
		`attachment; filename="`+file.OriginalName+`"`,
	)

	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "stream failed", http.StatusInternalServerError)
	}
}
