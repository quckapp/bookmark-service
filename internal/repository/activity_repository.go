package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(activity *model.BookmarkActivity) error
	GetByBookmark(bookmarkID uuid.UUID, limit, offset int) ([]model.BookmarkActivity, int64, error)
	GetByUser(userID uuid.UUID, limit, offset int) ([]model.BookmarkActivity, int64, error)
	GetRecent(userID uuid.UUID, limit int) ([]model.BookmarkActivity, error)
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	db.AutoMigrate(&model.BookmarkActivity{})
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(activity *model.BookmarkActivity) error {
	return r.db.Create(activity).Error
}

func (r *activityRepository) GetByBookmark(bookmarkID uuid.UUID, limit, offset int) ([]model.BookmarkActivity, int64, error) {
	var activities []model.BookmarkActivity
	var total int64

	r.db.Model(&model.BookmarkActivity{}).Where("bookmark_id = ?", bookmarkID).Count(&total)
	err := r.db.Where("bookmark_id = ?", bookmarkID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&activities).Error
	return activities, total, err
}

func (r *activityRepository) GetByUser(userID uuid.UUID, limit, offset int) ([]model.BookmarkActivity, int64, error) {
	var activities []model.BookmarkActivity
	var total int64

	r.db.Model(&model.BookmarkActivity{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&activities).Error
	return activities, total, err
}

func (r *activityRepository) GetRecent(userID uuid.UUID, limit int) ([]model.BookmarkActivity, error) {
	var activities []model.BookmarkActivity
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&activities).Error
	return activities, err
}
