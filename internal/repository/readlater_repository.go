package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type ReadLaterRepository interface {
	Create(item *model.ReadLaterItem) error
	GetByID(id uuid.UUID) (*model.ReadLaterItem, error)
	GetByUser(userID uuid.UUID, limit, offset int) ([]model.ReadLaterItem, int64, error)
	UpdateStatus(id uuid.UUID, status model.ReadLaterStatus) error
	Delete(id uuid.UUID) error
	GetStats(userID uuid.UUID) (*model.ReadLaterStats, error)
}

type readLaterRepository struct {
	db *gorm.DB
}

func NewReadLaterRepository(db *gorm.DB) ReadLaterRepository {
	db.AutoMigrate(&model.ReadLaterItem{})
	return &readLaterRepository{db: db}
}

func (r *readLaterRepository) Create(item *model.ReadLaterItem) error {
	return r.db.Create(item).Error
}

func (r *readLaterRepository) GetByID(id uuid.UUID) (*model.ReadLaterItem, error) {
	var item model.ReadLaterItem
	err := r.db.First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *readLaterRepository) GetByUser(userID uuid.UUID, limit, offset int) ([]model.ReadLaterItem, int64, error) {
	var items []model.ReadLaterItem
	var total int64

	r.db.Model(&model.ReadLaterItem{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("priority DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&items).Error
	return items, total, err
}

func (r *readLaterRepository) UpdateStatus(id uuid.UUID, status model.ReadLaterStatus) error {
	updates := map[string]interface{}{"status": status}
	if status == model.ReadLaterStatusCompleted {
		updates["read_at"] = gorm.Expr("NOW()")
	}
	return r.db.Model(&model.ReadLaterItem{}).Where("id = ?", id).Updates(updates).Error
}

func (r *readLaterRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ReadLaterItem{}, "id = ?", id).Error
}

func (r *readLaterRepository) GetStats(userID uuid.UUID) (*model.ReadLaterStats, error) {
	stats := &model.ReadLaterStats{}
	r.db.Model(&model.ReadLaterItem{}).Where("user_id = ?", userID).Count(&stats.Total)
	r.db.Model(&model.ReadLaterItem{}).Where("user_id = ? AND status = ?", userID, model.ReadLaterStatusUnread).Count(&stats.Unread)
	r.db.Model(&model.ReadLaterItem{}).Where("user_id = ? AND status = ?", userID, model.ReadLaterStatusReading).Count(&stats.Reading)
	r.db.Model(&model.ReadLaterItem{}).Where("user_id = ? AND status = ?", userID, model.ReadLaterStatusCompleted).Count(&stats.Completed)
	r.db.Model(&model.ReadLaterItem{}).Where("user_id = ? AND status = ?", userID, model.ReadLaterStatusArchived).Count(&stats.Archived)
	return stats, nil
}
