package image

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/model/dto"
)

func newTestImageService(t *testing.T) *ImageService {
	t.Helper()

	logger := zerolog.Nop()
	dir := t.TempDir()

	svc, err := New(&logger, dir)
	if err != nil {
		t.Fatalf("failed to create image service: %v", err)
	}

	return svc
}

func TestProcessImage_Base64Image(t *testing.T) {
	svc := newTestImageService(t)

	data := []byte("fake image data")
	encoded := base64.StdEncoding.EncodeToString(data)

	images := []dto.ImageData{
		{
			Base64: encoded,
			Type:   "image/png",
		},
	}

	paths, err := svc.ProcessImage(images)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(paths) != 1 {
		t.Fatalf("expected 1 image, got %d", len(paths))
	}

	if _, err := os.Stat(paths[0]); err != nil {
		t.Fatalf("image file not created: %v", err)
	}

	content, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != string(data) {
		t.Fatal("file content mismatch")
	}
}

func TestProcessImage_URLImage(t *testing.T) {
	svc := newTestImageService(t)

	images := []dto.ImageData{
		{
			URL: "https://example.com/image.jpg",
		},
	}

	paths, err := svc.ProcessImage(images)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if paths[0] != "https://example.com/image.jpg" {
		t.Fatalf("unexpected path: %s", paths[0])
	}
}

func TestProcessImage_FallbackOnFailure(t *testing.T) {
	svc := newTestImageService(t)

	validData := base64.StdEncoding.EncodeToString([]byte("ok"))
	invalidData := "%%%invalid-base64%%%"

	images := []dto.ImageData{
		{
			Base64: validData,
			Type:   "image/png",
		},
		{
			Base64: invalidData,
			Type:   "image/png",
		},
	}

	paths, err := svc.ProcessImage(images)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if paths != nil {
		t.Fatal("expected nil paths on failure")
	}

	files, err := os.ReadDir(svc.directory)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 0 {
		t.Fatal("rollback failed, files still exist")
	}
}

