package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type BookmarkReminderService interface {
	Create(reminder *model.BookmarkReminder) error
	GetByID(id uuid.UUID) (*model.BookmarkReminder, error)
	GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkReminder, error)
	GetByUser(userID uuid.UUID, page, limit int) ([]model.BookmarkReminder, int64, error)
	GetPending(userID uuid.UUID) ([]model.BookmarkReminder, error)
	Update(id uuid.UUID, remindAt time.Time, message string) (*model.BookmarkReminder, error)
	Delete(id uuid.UUID) error
	Cancel(id uuid.UUID) error
	CancelByBookmark(bookmarkID uuid.UUID) error
}

type bookmarkReminderService struct {
	repo   repository.ReminderRepository
	logger *zap.Logger
}

func NewBookmarkReminderService(repo repository.ReminderRepository, logger *zap.Logger) BookmarkReminderService {
	return &bookmarkReminderService{repo: repo, logger: logger}
}

func (s *bookmarkReminderService) Create(reminder *model.BookmarkReminder) error {
	if reminder.RemindAt.Before(time.Now()) {
		return fmt.Errorf("remind time must be in the future")
	}
	reminder.Status = "pending"
	err := s.repo.Create(reminder)
	if err != nil {
		return err
	}
	s.logger.Info("Created bookmark reminder", zap.String("id", reminder.ID.String()))
	return nil
}

func (s *bookmarkReminderService) GetByID(id uuid.UUID) (*model.BookmarkReminder, error) {
	return s.repo.GetByID(id)
}

func (s *bookmarkReminderService) GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkReminder, error) {
	return s.repo.GetByBookmark(bookmarkID)
}

func (s *bookmarkReminderService) GetByUser(userID uuid.UUID, page, limit int) ([]model.BookmarkReminder, int64, error) {
	offset := page * limit
	return s.repo.GetByUser(userID, limit, offset)
}

func (s *bookmarkReminderService) GetPending(userID uuid.UUID) ([]model.BookmarkReminder, error) {
	return s.repo.GetPending(userID)
}

func (s *bookmarkReminderService) Update(id uuid.UUID, remindAt time.Time, message string) (*model.BookmarkReminder, error) {
	reminder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("reminder not found")
	}
	if reminder.Status != "pending" {
		return nil, fmt.Errorf("can only update pending reminders")
	}

	reminder.RemindAt = remindAt
	if message != "" {
		reminder.Message = message
	}

	err = s.repo.Update(reminder)
	if err != nil {
		return nil, err
	}
	return reminder, nil
}

func (s *bookmarkReminderService) Delete(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	s.logger.Info("Deleted bookmark reminder", zap.String("id", id.String()))
	return nil
}

func (s *bookmarkReminderService) Cancel(id uuid.UUID) error {
	reminder, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("reminder not found")
	}
	if reminder.Status != "pending" {
		return fmt.Errorf("can only cancel pending reminders")
	}
	reminder.Status = "cancelled"
	return s.repo.Update(reminder)
}

func (s *bookmarkReminderService) CancelByBookmark(bookmarkID uuid.UUID) error {
	return s.repo.CancelByBookmark(bookmarkID)
}
