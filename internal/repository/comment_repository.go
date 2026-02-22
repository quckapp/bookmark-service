package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *model.BookmarkComment) error
	GetByID(id uuid.UUID) (*model.BookmarkComment, error)
	GetByBookmark(bookmarkID uuid.UUID, limit, offset int) ([]model.BookmarkComment, int64, error)
	Update(comment *model.BookmarkComment) error
	Delete(id uuid.UUID) error
	CountByBookmark(bookmarkID uuid.UUID) (int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	db.AutoMigrate(&model.BookmarkComment{})
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *model.BookmarkComment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByID(id uuid.UUID) (*model.BookmarkComment, error) {
	var comment model.BookmarkComment
	err := r.db.First(&comment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) GetByBookmark(bookmarkID uuid.UUID, limit, offset int) ([]model.BookmarkComment, int64, error) {
	var comments []model.BookmarkComment
	var total int64

	r.db.Model(&model.BookmarkComment{}).Where("bookmark_id = ?", bookmarkID).Count(&total)
	err := r.db.Where("bookmark_id = ?", bookmarkID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&comments).Error
	return comments, total, err
}

func (r *commentRepository) Update(comment *model.BookmarkComment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uuid.UUID) error {
	r.db.Where("parent_id = ?", id).Delete(&model.BookmarkComment{})
	return r.db.Delete(&model.BookmarkComment{}, "id = ?", id).Error
}

func (r *commentRepository) CountByBookmark(bookmarkID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.BookmarkComment{}).Where("bookmark_id = ?", bookmarkID).Count(&count).Error
	return count, err
}
