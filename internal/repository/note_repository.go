package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type NoteRepository interface {
	Create(note *model.BookmarkNote) error
	GetByID(id uuid.UUID) (*model.BookmarkNote, error)
	GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error)
	GetByUser(userID uuid.UUID, limit, offset int) ([]model.BookmarkNote, int64, error)
	Update(note *model.BookmarkNote) error
	Delete(id uuid.UUID) error
	GetPinnedByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error)
	CountByBookmark(bookmarkID uuid.UUID) (int64, error)
}

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	db.AutoMigrate(&model.BookmarkNote{})
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(note *model.BookmarkNote) error {
	return r.db.Create(note).Error
}

func (r *noteRepository) GetByID(id uuid.UUID) (*model.BookmarkNote, error) {
	var note model.BookmarkNote
	err := r.db.First(&note, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) GetByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error) {
	var notes []model.BookmarkNote
	err := r.db.Where("bookmark_id = ?", bookmarkID).
		Order("is_pinned DESC, created_at DESC").
		Find(&notes).Error
	return notes, err
}

func (r *noteRepository) GetByUser(userID uuid.UUID, limit, offset int) ([]model.BookmarkNote, int64, error) {
	var notes []model.BookmarkNote
	var total int64

	r.db.Model(&model.BookmarkNote{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&notes).Error
	return notes, total, err
}

func (r *noteRepository) Update(note *model.BookmarkNote) error {
	return r.db.Save(note).Error
}

func (r *noteRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.BookmarkNote{}, "id = ?", id).Error
}

func (r *noteRepository) GetPinnedByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkNote, error) {
	var notes []model.BookmarkNote
	err := r.db.Where("bookmark_id = ? AND is_pinned = ?", bookmarkID, true).
		Order("created_at DESC").
		Find(&notes).Error
	return notes, err
}

func (r *noteRepository) CountByBookmark(bookmarkID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.BookmarkNote{}).Where("bookmark_id = ?", bookmarkID).Count(&count).Error
	return count, err
}
