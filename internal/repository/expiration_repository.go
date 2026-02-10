package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type ExpirationRepository interface {
	Create(expiration *model.BookmarkExpiration) error
	GetByBookmarkID(bookmarkID uuid.UUID) (*model.BookmarkExpiration, error)
	Delete(bookmarkID uuid.UUID) error
	GetExpiring(userID uuid.UUID, before time.Time) ([]model.BookmarkExpiration, error)
	MarkExpired(id uuid.UUID) error
	Update(expiration *model.BookmarkExpiration) error
}

type expirationRepository struct {
	db *gorm.DB
}

func NewExpirationRepository(db *gorm.DB) ExpirationRepository {
	db.AutoMigrate(&model.BookmarkExpiration{})
	return &expirationRepository{db: db}
}

func (r *expirationRepository) Create(expiration *model.BookmarkExpiration) error {
	return r.db.Create(expiration).Error
}

func (r *expirationRepository) GetByBookmarkID(bookmarkID uuid.UUID) (*model.BookmarkExpiration, error) {
	var expiration model.BookmarkExpiration
	err := r.db.Where("bookmark_id = ?", bookmarkID).First(&expiration).Error
	if err != nil {
		return nil, err
	}
	return &expiration, nil
}

func (r *expirationRepository) Delete(bookmarkID uuid.UUID) error {
	return r.db.Where("bookmark_id = ?", bookmarkID).Delete(&model.BookmarkExpiration{}).Error
}

func (r *expirationRepository) GetExpiring(userID uuid.UUID, before time.Time) ([]model.BookmarkExpiration, error) {
	var expirations []model.BookmarkExpiration
	err := r.db.Where("user_id = ? AND expires_at <= ? AND is_expired = ?", userID, before, false).
		Order("expires_at ASC").
		Find(&expirations).Error
	return expirations, err
}

func (r *expirationRepository) MarkExpired(id uuid.UUID) error {
	return r.db.Model(&model.BookmarkExpiration{}).Where("id = ?", id).Update("is_expired", true).Error
}

func (r *expirationRepository) Update(expiration *model.BookmarkExpiration) error {
	return r.db.Save(expiration).Error
}
