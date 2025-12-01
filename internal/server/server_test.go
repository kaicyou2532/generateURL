package server

import (
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
