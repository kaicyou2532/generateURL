package local

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestSaverSave(t *testing.T) {
	dir := t.TempDir()

	saver, err := New(dir)
	if err != nil {
		t.Fatalf("New saver: %v", err)
	}

	src := filepath.Join(dir, "input.bin")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	file, err := os.Open(src)
	if err != nil {
		t.Fatalf("open source: %v", err)
	}

	header := &multipart.FileHeader{Filename: "picture.jpg"}

	name, err := saver.Save(context.Background(), file, header)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	if filepath.Ext(name) != ".jpg" {
		t.Fatalf("expected .jpg extension, got %s", name)
	}

	target := filepath.Join(dir, name)
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("stat saved file: %v", err)
	}
}
