package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type ReminderRepository interface {
	Create(reminder *model.BookmarkReminder) error
	GetByID(id uuid.UUID) (*model.BookmarkReminder, error)
	GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkReminder, error)
	GetByUser(userID uuid.UUID, limit, offset int) ([]model.BookmarkReminder, int64, error)
	GetPending(userID uuid.UUID) ([]model.BookmarkReminder, error)
	Update(reminder *model.BookmarkReminder) error
	Delete(id uuid.UUID) error
	MarkFired(id uuid.UUID) error
	GetDueBefore(t time.Time) ([]model.BookmarkReminder, error)
	CancelByBookmark(bookmarkID uuid.UUID) error
}

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	db.AutoMigrate(&model.BookmarkReminder{})
	return &reminderRepository{db: db}
}

func (r *reminderRepository) Create(reminder *model.BookmarkReminder) error {
	return r.db.Create(reminder).Error
}

func (r *reminderRepository) GetByID(id uuid.UUID) (*model.BookmarkReminder, error) {
	var reminder model.BookmarkReminder
	err := r.db.First(&reminder, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkReminder, error) {
	var reminders []model.BookmarkReminder
	err := r.db.Where("bookmark_id = ?", bookmarkID).Order("remind_at ASC").Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetByUser(userID uuid.UUID, limit, offset int) ([]model.BookmarkReminder, int64, error) {
	var reminders []model.BookmarkReminder
	var total int64

	r.db.Model(&model.BookmarkReminder{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("remind_at ASC").
		Limit(limit).Offset(offset).
		Find(&reminders).Error
	return reminders, total, err
}

func (r *reminderRepository) GetPending(userID uuid.UUID) ([]model.BookmarkReminder, error) {
	var reminders []model.BookmarkReminder
	err := r.db.Where("user_id = ? AND status = ?", userID, "pending").
		Order("remind_at ASC").
		Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) Update(reminder *model.BookmarkReminder) error {
	return r.db.Save(reminder).Error
}

func (r *reminderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.BookmarkReminder{}, "id = ?", id).Error
}

func (r *reminderRepository) MarkFired(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.BookmarkReminder{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "fired", "fired_at": now}).Error
}

func (r *reminderRepository) GetDueBefore(t time.Time) ([]model.BookmarkReminder, error) {
	var reminders []model.BookmarkReminder
	err := r.db.Where("status = ? AND remind_at <= ?", "pending", t).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) CancelByBookmark(bookmarkID uuid.UUID) error {
	return r.db.Model(&model.BookmarkReminder{}).
		Where("bookmark_id = ? AND status = ?", bookmarkID, "pending").
		Update("status", "cancelled").Error
}
