package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type FavoriteService interface {
	AddFavorite(userID, bookmarkID uuid.UUID) error
	RemoveFavorite(userID, bookmarkID uuid.UUID) error
	IsFavorite(userID, bookmarkID uuid.UUID) (bool, error)
	GetFavorites(userID uuid.UUID, page, limit int) ([]model.Bookmark, int64, error)
	GetFavoriteCount(userID uuid.UUID) (int64, error)
}

type favoriteService struct {
	repo   repository.FavoriteRepository
	logger *zap.Logger
}

func NewFavoriteService(repo repository.FavoriteRepository, logger *zap.Logger) FavoriteService {
	return &favoriteService{repo: repo, logger: logger}
}

func (s *favoriteService) AddFavorite(userID, bookmarkID uuid.UUID) error {
	isFav, _ := s.repo.IsFavorite(userID, bookmarkID)
	if isFav {
		return fmt.Errorf("already favorited")
	}
	fav := &model.BookmarkFavorite{
		UserID:     userID,
		BookmarkID: bookmarkID,
	}
	return s.repo.Create(fav)
}

func (s *favoriteService) RemoveFavorite(userID, bookmarkID uuid.UUID) error {
	return s.repo.Delete(userID, bookmarkID)
}

func (s *favoriteService) IsFavorite(userID, bookmarkID uuid.UUID) (bool, error) {
	return s.repo.IsFavorite(userID, bookmarkID)
}

func (s *favoriteService) GetFavorites(userID uuid.UUID, page, limit int) ([]model.Bookmark, int64, error) {
	offset := page * limit
	return s.repo.GetByUser(userID, limit, offset)
}

func (s *favoriteService) GetFavoriteCount(userID uuid.UUID) (int64, error) {
	return s.repo.Count(userID)
}
