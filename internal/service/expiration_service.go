package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type ExpirationService interface {
	Set(expiration *model.BookmarkExpiration) error
	GetByBookmarkID(bookmarkID uuid.UUID) (*model.BookmarkExpiration, error)
	Remove(bookmarkID uuid.UUID) error
	GetExpiring(userID uuid.UUID) ([]model.BookmarkExpiration, error)
}

type expirationService struct {
	repo   repository.ExpirationRepository
	logger *zap.Logger
}

func NewExpirationService(repo repository.ExpirationRepository, logger *zap.Logger) ExpirationService {
	return &expirationService{repo: repo, logger: logger}
}

func (s *expirationService) Set(expiration *model.BookmarkExpiration) error {
	// Check if expiration already exists for this bookmark
	existing, _ := s.repo.GetByBookmarkID(expiration.BookmarkID)
	if existing != nil {
		existing.ExpiresAt = expiration.ExpiresAt
		existing.Action = expiration.Action
		existing.IsExpired = false
		return s.repo.Update(existing)
	}

	err := s.repo.Create(expiration)
	if err != nil {
		return err
	}
	s.logger.Info("Set expiration", zap.String("bookmarkId", expiration.BookmarkID.String()))
	return nil
}

func (s *expirationService) GetByBookmarkID(bookmarkID uuid.UUID) (*model.BookmarkExpiration, error) {
	return s.repo.GetByBookmarkID(bookmarkID)
}

func (s *expirationService) Remove(bookmarkID uuid.UUID) error {
	return s.repo.Delete(bookmarkID)
}

func (s *expirationService) GetExpiring(userID uuid.UUID) ([]model.BookmarkExpiration, error) {
	// Get items expiring in next 7 days
	before := time.Now().AddDate(0, 0, 7)
	return s.repo.GetExpiring(userID, before)
}
