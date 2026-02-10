package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type FavoriteRepository interface {
	Create(fav *model.BookmarkFavorite) error
	Delete(userID, bookmarkID uuid.UUID) error
	IsFavorite(userID, bookmarkID uuid.UUID) (bool, error)
	GetByUser(userID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error)
	Count(userID uuid.UUID) (int64, error)
}

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	db.AutoMigrate(&model.BookmarkFavorite{})
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Create(fav *model.BookmarkFavorite) error {
	return r.db.Create(fav).Error
}

func (r *favoriteRepository) Delete(userID, bookmarkID uuid.UUID) error {
	return r.db.Where("user_id = ? AND bookmark_id = ?", userID, bookmarkID).
		Delete(&model.BookmarkFavorite{}).Error
}

func (r *favoriteRepository) IsFavorite(userID, bookmarkID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.BookmarkFavorite{}).
		Where("user_id = ? AND bookmark_id = ?", userID, bookmarkID).
		Count(&count).Error
	return count > 0, err
}

func (r *favoriteRepository) GetByUser(userID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error) {
	var bookmarks []model.Bookmark
	var total int64

	r.db.Model(&model.BookmarkFavorite{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Joins("JOIN bookmark_favorites ON bookmark_favorites.bookmark_id = bookmarks.id").
		Where("bookmark_favorites.user_id = ?", userID).
		Order("bookmark_favorites.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&bookmarks).Error

	return bookmarks, total, err
}

func (r *favoriteRepository) Count(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.BookmarkFavorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
