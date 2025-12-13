package image

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/model/dto"
)

type ImageService struct {
	logger    *zerolog.Logger
	directory string
}

func New(l *zerolog.Logger) (*ImageService, error) {
	dir := "../uploads" // can be more dynamic
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	return &ImageService{
		logger:    l,
		directory: dir,
	}, nil
}

func (s *ImageService) ProcessImage(images []dto.ImageData) ([]string, error) {
	// can add wg, chan and gr to make fast.
	processedImage := make([]string, len(images))
	ifFailedImage := make([]string, len(images))
	for _, img := range images {
		if img.URL != "" {
			processedImage = append(processedImage, img.URL)
			continue
		}

		urlPath, err := s.storeImage(img)
		if err != nil {
			// romove all uploded imge
			s.fallback(ifFailedImage)
			return nil, err
		}
		ifFailedImage = append(ifFailedImage, urlPath)
		processedImage = append(processedImage, urlPath)
	}

	s.logger.Info().Str("event", "image upload").Int("total", len(ifFailedImage)).Msg("upload complate")
	return processedImage, nil

}

func (s *ImageService) storeImage(img dto.ImageData) (string, error) {
	data, err := base64.StdEncoding.DecodeString(img.Base64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	hash := md5.Sum(data)
	filename := hex.EncodeToString(hash[:]) + s.getExtension(img.Type)
	filepath := filepath.Join(s.directory, filename)

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return filepath, nil
}

func (s *ImageService) fallback(img []string) {
	for _, f := range img {
		_ = os.Remove(f)
	}
}

func (is *ImageService) getExtension(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

