package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kaicyou2532/generateURL/internal/config"
	"github.com/kaicyou2532/generateURL/internal/storage"
)

// Server exposes HTTP endpoints for uploading images and serving them back.
type Server struct {
	cfg   config.Config
	saver storage.Saver
	mux   *http.ServeMux
}

// New constructs a Server with the provided dependencies.
func New(cfg config.Config, saver storage.Saver) (*Server, error) {
	if saver == nil {
		return nil, errors.New("storage saver is required")
	}

	s := &Server{cfg: cfg, saver: saver, mux: http.NewServeMux()}
	s.registerRoutes()
	return s, nil
}

// Handler returns the configured HTTP handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	fileServer := http.StripPrefix("/files/", http.FileServer(http.Dir(s.cfg.UploadDir)))
	s.mux.Handle("GET /files/", fileServer)

	s.mux.HandleFunc("POST /api/v1/uploads", s.withCORS(s.handleUpload))
	s.mux.HandleFunc("OPTIONS /api/v1/uploads", s.withCORS(func(http.ResponseWriter, *http.Request) {}))
}

func (s *Server) withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.cfg.MaxFileSize)
	if err := r.ParseMultipartForm(s.cfg.MaxFileSize); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("parse form: %w", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("file form field missing: %w", err))
		return
	}

	name, err := s.saver.Save(r.Context(), file, header)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	url := s.buildFileURL(r, name)
	writeJSON(w, http.StatusCreated, map[string]string{"url": url})
}

func (s *Server) buildFileURL(r *http.Request, fileName string) string {
	base := s.cfg.BaseURL
	if base == "" {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		base = fmt.Sprintf("%s://%s", scheme, r.Host)
	}

	return fmt.Sprintf("%s/files/%s", strings.TrimSuffix(base, "/"), fileName)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
