package service

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type CommentService interface {
	Create(comment *model.BookmarkComment) error
	GetByBookmark(bookmarkID uuid.UUID, page, limit int) ([]model.BookmarkComment, int64, error)
	Update(id uuid.UUID, content string) (*model.BookmarkComment, error)
	Delete(id uuid.UUID) error
}

type commentService struct {
	repo   repository.CommentRepository
	logger *zap.Logger
}

func NewCommentService(repo repository.CommentRepository, logger *zap.Logger) CommentService {
	return &commentService{repo: repo, logger: logger}
}

func (s *commentService) Create(comment *model.BookmarkComment) error {
	err := s.repo.Create(comment)
	if err != nil {
		return err
	}
	s.logger.Info("Created comment", zap.String("bookmarkId", comment.BookmarkID.String()))
	return nil
}

func (s *commentService) GetByBookmark(bookmarkID uuid.UUID, page, limit int) ([]model.BookmarkComment, int64, error) {
	offset := page * limit
	return s.repo.GetByBookmark(bookmarkID, limit, offset)
}

func (s *commentService) Update(id uuid.UUID, content string) (*model.BookmarkComment, error) {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	comment.Content = content
	err = s.repo.Update(comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *commentService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
