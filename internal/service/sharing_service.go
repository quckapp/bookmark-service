package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type SharingService interface {
	ShareBookmark(shared *model.SharedBookmark) error
	GetSharedWithUser(userID uuid.UUID, page, limit int) ([]model.SharedBookmark, int64, error)
	GetSharedByUser(userID uuid.UUID, page, limit int) ([]model.SharedBookmark, int64, error)
	AcceptShare(id uuid.UUID) error
	DeclineShare(id uuid.UUID) error
	GetPendingCount(userID uuid.UUID) (int64, error)
}

type sharingService struct {
	repo         repository.SharingRepository
	bookmarkRepo repository.BookmarkRepository
	logger       *zap.Logger
}

func NewSharingService(repo repository.SharingRepository, bookmarkRepo repository.BookmarkRepository, logger *zap.Logger) SharingService {
	return &sharingService{repo: repo, bookmarkRepo: bookmarkRepo, logger: logger}
}

func (s *sharingService) ShareBookmark(shared *model.SharedBookmark) error {
	// Verify bookmark exists
	_, err := s.bookmarkRepo.GetByID(shared.BookmarkID)
	if err != nil {
		return fmt.Errorf("bookmark not found")
	}

	// Check if already shared
	alreadyShared, _ := s.repo.IsAlreadyShared(shared.BookmarkID, shared.SharedWith)
	if alreadyShared {
		return fmt.Errorf("bookmark already shared with this user")
	}

	err = s.repo.Create(shared)
	if err != nil {
		return err
	}
	s.logger.Info("Shared bookmark",
		zap.String("bookmarkID", shared.BookmarkID.String()),
		zap.String("sharedWith", shared.SharedWith.String()))
	return nil
}

func (s *sharingService) GetSharedWithUser(userID uuid.UUID, page, limit int) ([]model.SharedBookmark, int64, error) {
	offset := page * limit
	return s.repo.GetSharedWithUser(userID, limit, offset)
}

func (s *sharingService) GetSharedByUser(userID uuid.UUID, page, limit int) ([]model.SharedBookmark, int64, error) {
	offset := page * limit
	return s.repo.GetSharedByUser(userID, limit, offset)
}

func (s *sharingService) AcceptShare(id uuid.UUID) error {
	share, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("shared bookmark not found")
	}
	if share.IsAccepted {
		return fmt.Errorf("already accepted")
	}
	return s.repo.Accept(id)
}

func (s *sharingService) DeclineShare(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("shared bookmark not found")
	}
	return s.repo.Delete(id)
}

func (s *sharingService) GetPendingCount(userID uuid.UUID) (int64, error) {
	return s.repo.GetPendingCount(userID)
}
