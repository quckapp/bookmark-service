package service

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type ReadLaterService interface {
	Add(item *model.ReadLaterItem) error
	GetByUser(userID uuid.UUID, page, limit int) ([]model.ReadLaterItem, int64, error)
	UpdateStatus(id uuid.UUID, status model.ReadLaterStatus) error
	Delete(id uuid.UUID) error
	GetStats(userID uuid.UUID) (*model.ReadLaterStats, error)
}

type readLaterService struct {
	repo   repository.ReadLaterRepository
	logger *zap.Logger
}

func NewReadLaterService(repo repository.ReadLaterRepository, logger *zap.Logger) ReadLaterService {
	return &readLaterService{repo: repo, logger: logger}
}

func (s *readLaterService) Add(item *model.ReadLaterItem) error {
	if item.Status == "" {
		item.Status = model.ReadLaterStatusUnread
	}
	err := s.repo.Create(item)
	if err != nil {
		return err
	}
	s.logger.Info("Added to read later", zap.String("bookmarkId", item.BookmarkID.String()))
	return nil
}

func (s *readLaterService) GetByUser(userID uuid.UUID, page, limit int) ([]model.ReadLaterItem, int64, error) {
	offset := page * limit
	return s.repo.GetByUser(userID, limit, offset)
}

func (s *readLaterService) UpdateStatus(id uuid.UUID, status model.ReadLaterStatus) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *readLaterService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *readLaterService) GetStats(userID uuid.UUID) (*model.ReadLaterStats, error) {
	return s.repo.GetStats(userID)
}
