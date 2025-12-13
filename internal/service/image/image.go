package image

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/model/dto"
)

type ImageService struct {
	logger *zerolog.Logger
}

func New(l *zerolog.Logger) *ImageService {
	return &ImageService{
		logger: l,
	}
}

func (s *ImageService) ProcessImage(images []dto.ImageData) []string {
	processedImage := make([]string, 0, len(images))
	counter := 0

	for _, img := range images {
		if img.URL != "" {
			processedImage = append(processedImage, img.URL)
		}
		counter++
		urlPath := s.storeImage(img, counter)
		processedImage = append(processedImage, urlPath)
	}

	s.logger.Info().Str("event", "image upload").Msg("upload complate")
	return processedImage
}

func (s *ImageService) storeImage(img dto.ImageData, counter int) string {
	fullname := fmt.Sprintf("localImg%d%s", counter, s.getExtension(img.Type))
	return fullname
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
