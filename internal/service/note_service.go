package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type NoteService interface {
	Create(note *model.BookmarkNote) error
	GetByID(id uuid.UUID) (*model.BookmarkNote, error)
	GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error)
	GetByUser(userID uuid.UUID, page, limit int) ([]model.BookmarkNote, int64, error)
	Update(id uuid.UUID, req *model.UpdateNoteRequest) (*model.BookmarkNote, error)
	Delete(id uuid.UUID) error
	GetPinnedByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error)
}

type noteService struct {
	repo   repository.NoteRepository
	logger *zap.Logger
}

func NewNoteService(repo repository.NoteRepository, logger *zap.Logger) NoteService {
	return &noteService{repo: repo, logger: logger}
}

func (s *noteService) Create(note *model.BookmarkNote) error {
	err := s.repo.Create(note)
	if err != nil {
		return err
	}
	s.logger.Info("Created note", zap.String("id", note.ID.String()))
	return nil
}

func (s *noteService) GetByID(id uuid.UUID) (*model.BookmarkNote, error) {
	return s.repo.GetByID(id)
}

func (s *noteService) GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error) {
	return s.repo.GetByBookmark(bookmarkID)
}

func (s *noteService) GetByUser(userID uuid.UUID, page, limit int) ([]model.BookmarkNote, int64, error) {
	offset := page * limit
	return s.repo.GetByUser(userID, limit, offset)
}

func (s *noteService) Update(id uuid.UUID, req *model.UpdateNoteRequest) (*model.BookmarkNote, error) {
	note, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("note not found")
	}

	if req.Content != "" {
		note.Content = req.Content
	}
	if req.Color != "" {
		note.Color = req.Color
	}
	if req.IsPinned != nil {
		note.IsPinned = *req.IsPinned
	}

	err = s.repo.Update(note)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (s *noteService) Delete(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	s.logger.Info("Deleted note", zap.String("id", id.String()))
	return nil
}

func (s *noteService) GetPinnedByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error) {
	return s.repo.GetPinnedByBookmark(bookmarkID)
}
