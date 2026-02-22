package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type PreviewService interface {
	Create(preview *model.LinkPreview) error
	GetByBookmarkID(bookmarkID uuid.UUID) (*model.LinkPreview, error)
	GetByURL(url string) (*model.LinkPreview, error)
	Generate(bookmarkID uuid.UUID, url string) (*model.LinkPreview, error)
}

type previewService struct {
	repo repository.PreviewRepository
	logger *zap.Logger
}

func NewPreviewService(repo repository.PreviewRepository, logger *zap.Logger) PreviewService {
	return &previewService{repo: repo, logger: logger}
}

func (s *previewService) Create(preview *model.LinkPreview) error {
	preview.FetchedAt = time.Now()
	return s.repo.Create(preview)
}

func (s *previewService) GetByBookmarkID(bookmarkID uuid.UUID) (*model.LinkPreview, error) {
	return s.repo.GetByBookmarkID(bookmarkID)
}

func (s *previewService) GetByURL(url string) (*model.LinkPreview, error) {
	return s.repo.GetByURL(url)
}

func (s *previewService) Generate(bookmarkID uuid.UUID, url string) (*model.LinkPreview, error) {
	// Check if preview already exists
	existing, _ := s.repo.GetByBookmarkID(bookmarkID)
	if existing != nil {
		existing.URL = url
		existing.FetchedAt = time.Now()
		err := s.repo.Update(existing)
		return existing, err
	}

	preview := &model.LinkPreview{
		BookmarkID: bookmarkID,
		URL:        url,
		FetchedAt:  time.Now(),
	}
	err := s.repo.Create(preview)
	if err != nil {
		return nil, err
	}
	s.logger.Info("Generated preview", zap.String("bookmarkId", bookmarkID.String()))
	return preview, nil
}
