package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type PreviewRepository interface {
	Create(preview *model.LinkPreview) error
	GetByBookmarkID(bookmarkID uuid.UUID) (*model.LinkPreview, error)
	GetByURL(url string) (*model.LinkPreview, error)
	Update(preview *model.LinkPreview) error
	Delete(bookmarkID uuid.UUID) error
}

type previewRepository struct {
	db *gorm.DB
}

func NewPreviewRepository(db *gorm.DB) PreviewRepository {
	db.AutoMigrate(&model.LinkPreview{})
	return &previewRepository{db: db}
}

func (r *previewRepository) Create(preview *model.LinkPreview) error {
	return r.db.Create(preview).Error
}

func (r *previewRepository) GetByBookmarkID(bookmarkID uuid.UUID) (*model.LinkPreview, error) {
	var preview model.LinkPreview
	err := r.db.Where("bookmark_id = ?", bookmarkID).First(&preview).Error
	if err != nil {
		return nil, err
	}
	return &preview, nil
}

func (r *previewRepository) GetByURL(url string) (*model.LinkPreview, error) {
	var preview model.LinkPreview
	err := r.db.Where("url = ?", url).First(&preview).Error
	if err != nil {
		return nil, err
	}
	return &preview, nil
}

func (r *previewRepository) Update(preview *model.LinkPreview) error {
	return r.db.Save(preview).Error
}

func (r *previewRepository) Delete(bookmarkID uuid.UUID) error {
	return r.db.Where("bookmark_id = ?", bookmarkID).Delete(&model.LinkPreview{}).Error
}
