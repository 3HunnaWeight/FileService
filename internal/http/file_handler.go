package http

import (
	"context"
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
	Delete(ctx context.Context, publicID string) error
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
		respondErr(w, http.StatusBadRequest, "missing file field")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		respondErr(w, http.StatusBadRequest, "failed to read file")
		return
	}

	publicID, err := h.service.Upload(
		r.Context(),
		header.Filename,
		data,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, "upload failed")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"id":  publicID,
		"url": "/files/" + publicID,
	})
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	publicID := chi.URLParam(r, "id")

	reader, file, err := h.service.Download(r.Context(), publicID)
	if err != nil {
		respondErr(w, http.StatusNotFound, "file not found")
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+file.OriginalName+`"`)

	if _, err := io.Copy(w, reader); err != nil {
		respondErr(w, http.StatusInternalServerError, "stream failed")
	}
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	publicID := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), publicID); err != nil {
		respondErr(w, http.StatusNotFound, "file not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
