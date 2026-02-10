package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type VersionRepository interface {
	Create(version *model.BookmarkVersion) error
	GetByID(id uuid.UUID) (*model.BookmarkVersion, error)
	GetByBookmark(bookmarkID uuid.UUID, limit, offset int) ([]model.BookmarkVersion, int64, error)
	GetLatestVersion(bookmarkID uuid.UUID) (int, error)
}

type versionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) VersionRepository {
	db.AutoMigrate(&model.BookmarkVersion{})
	return &versionRepository{db: db}
}

func (r *versionRepository) Create(version *model.BookmarkVersion) error {
	return r.db.Create(version).Error
}

func (r *versionRepository) GetByID(id uuid.UUID) (*model.BookmarkVersion, error) {
	var version model.BookmarkVersion
	err := r.db.First(&version, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *versionRepository) GetByBookmark(bookmarkID uuid.UUID, limit, offset int) ([]model.BookmarkVersion, int64, error) {
	var versions []model.BookmarkVersion
	var total int64

	r.db.Model(&model.BookmarkVersion{}).Where("bookmark_id = ?", bookmarkID).Count(&total)
	err := r.db.Where("bookmark_id = ?", bookmarkID).
		Order("version DESC").
		Limit(limit).Offset(offset).
		Find(&versions).Error
	return versions, total, err
}

func (r *versionRepository) GetLatestVersion(bookmarkID uuid.UUID) (int, error) {
	var version model.BookmarkVersion
	err := r.db.Where("bookmark_id = ?", bookmarkID).Order("version DESC").First(&version).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return version.Version, nil
}
