package server

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaicyou2532/generateURL/internal/config"
)

func TestBuildFileURL_DefaultBase(t *testing.T) {
	cfg := config.Config{UploadDir: t.TempDir()}
	s := &Server{cfg: cfg}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Host = "example.com"

	url := s.buildFileURL(r, "abc.png")
	expected := "http://example.com/files/abc.png"
	if url != expected {
		t.Fatalf("expected %s, got %s", expected, url)
	}
}

func TestBuildFileURL_CustomBase(t *testing.T) {
	cfg := config.Config{UploadDir: t.TempDir(), BaseURL: "https://cdn.example.com"}
	s := &Server{cfg: cfg}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	url := s.buildFileURL(r, "abc.png")
	expected := "https://cdn.example.com/files/abc.png"
	if url != expected {
		t.Fatalf("expected %s, got %s", expected, url)
	}
}

func TestUploadsCORSPreflight(t *testing.T) {
	cfg := config.Config{UploadDir: t.TempDir()}
	s, err := New(cfg, stubSaver{})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/uploads", nil)
	rec := httptest.NewRecorder()

	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Fatalf("unexpected origin header: %q", origin)
	}
}

type stubSaver struct{}

func (stubSaver) Save(context.Context, multipart.File, *multipart.FileHeader) (string, error) {
	return "stub.png", nil
}
