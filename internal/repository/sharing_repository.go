package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type SharingRepository interface {
	Create(shared *model.SharedBookmark) error
	GetByID(id uuid.UUID) (*model.SharedBookmark, error)
	GetSharedWithUser(userID uuid.UUID, limit, offset int) ([]model.SharedBookmark, int64, error)
	GetSharedByUser(userID uuid.UUID, limit, offset int) ([]model.SharedBookmark, int64, error)
	Accept(id uuid.UUID) error
	Delete(id uuid.UUID) error
	IsAlreadyShared(bookmarkID, sharedWith uuid.UUID) (bool, error)
	GetPendingCount(userID uuid.UUID) (int64, error)
}

type sharingRepository struct {
	db *gorm.DB
}

func NewSharingRepository(db *gorm.DB) SharingRepository {
	db.AutoMigrate(&model.SharedBookmark{})
	return &sharingRepository{db: db}
}

func (r *sharingRepository) Create(shared *model.SharedBookmark) error {
	return r.db.Create(shared).Error
}

func (r *sharingRepository) GetByID(id uuid.UUID) (*model.SharedBookmark, error) {
	var shared model.SharedBookmark
	err := r.db.First(&shared, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &shared, nil
}

func (r *sharingRepository) GetSharedWithUser(userID uuid.UUID, limit, offset int) ([]model.SharedBookmark, int64, error) {
	var shares []model.SharedBookmark
	var total int64

	r.db.Model(&model.SharedBookmark{}).Where("shared_with = ?", userID).Count(&total)
	err := r.db.Where("shared_with = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&shares).Error
	return shares, total, err
}

func (r *sharingRepository) GetSharedByUser(userID uuid.UUID, limit, offset int) ([]model.SharedBookmark, int64, error) {
	var shares []model.SharedBookmark
	var total int64

	r.db.Model(&model.SharedBookmark{}).Where("shared_by = ?", userID).Count(&total)
	err := r.db.Where("shared_by = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&shares).Error
	return shares, total, err
}

func (r *sharingRepository) Accept(id uuid.UUID) error {
	return r.db.Model(&model.SharedBookmark{}).Where("id = ?", id).Update("is_accepted", true).Error
}

func (r *sharingRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.SharedBookmark{}, "id = ?", id).Error
}

func (r *sharingRepository) IsAlreadyShared(bookmarkID, sharedWith uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.SharedBookmark{}).
		Where("bookmark_id = ? AND shared_with = ?", bookmarkID, sharedWith).
		Count(&count).Error
	return count > 0, err
}

func (r *sharingRepository) GetPendingCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.SharedBookmark{}).
		Where("shared_with = ? AND is_accepted = ?", userID, false).
		Count(&count).Error
	return count, err
}
