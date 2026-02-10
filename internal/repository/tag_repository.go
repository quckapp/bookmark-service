package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(tag *model.BookmarkTag) error
	GetByID(id uuid.UUID) (*model.BookmarkTag, error)
	GetByUser(userID uuid.UUID) ([]model.BookmarkTag, error)
	GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTag, error)
	Update(tag *model.BookmarkTag) error
	Delete(id uuid.UUID) error
	AddTagToBookmark(mapping *model.BookmarkTagMapping) error
	RemoveTagFromBookmark(bookmarkID, tagID uuid.UUID) error
	GetTagsByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkTag, error)
	GetBookmarksByTag(tagID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error)
	RemoveAllTagsFromBookmark(bookmarkID uuid.UUID) error
	BulkAddTagToBookmarks(bookmarkIDs []uuid.UUID, tagID uuid.UUID) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	db.AutoMigrate(&model.BookmarkTag{}, &model.BookmarkTagMapping{})
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *model.BookmarkTag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) GetByID(id uuid.UUID) (*model.BookmarkTag, error) {
	var tag model.BookmarkTag
	err := r.db.First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetByUser(userID uuid.UUID) ([]model.BookmarkTag, error) {
	var tags []model.BookmarkTag
	err := r.db.Where("user_id = ?", userID).Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTag, error) {
	var tags []model.BookmarkTag
	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) Update(tag *model.BookmarkTag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uuid.UUID) error {
	// Remove all mappings first
	r.db.Where("tag_id = ?", id).Delete(&model.BookmarkTagMapping{})
	return r.db.Delete(&model.BookmarkTag{}, "id = ?", id).Error
}

func (r *tagRepository) AddTagToBookmark(mapping *model.BookmarkTagMapping) error {
	return r.db.Create(mapping).Error
}

func (r *tagRepository) RemoveTagFromBookmark(bookmarkID, tagID uuid.UUID) error {
	return r.db.Where("bookmark_id = ? AND tag_id = ?", bookmarkID, tagID).Delete(&model.BookmarkTagMapping{}).Error
}

func (r *tagRepository) GetTagsByBookmark(bookmarkID uuid.UUID) ([]model.BookmarkTag, error) {
	var tags []model.BookmarkTag
	err := r.db.Joins("JOIN bookmark_tag_mappings ON bookmark_tag_mappings.tag_id = bookmark_tags.id").
		Where("bookmark_tag_mappings.bookmark_id = ?", bookmarkID).
		Find(&tags).Error
	return tags, err
}

func (r *tagRepository) GetBookmarksByTag(tagID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error) {
	var bookmarks []model.Bookmark
	var total int64

	r.db.Model(&model.BookmarkTagMapping{}).Where("tag_id = ?", tagID).Count(&total)
	err := r.db.Joins("JOIN bookmark_tag_mappings ON bookmark_tag_mappings.bookmark_id = bookmarks.id").
		Where("bookmark_tag_mappings.tag_id = ?", tagID).
		Order("bookmarks.created_at DESC").
		Limit(limit).Offset(offset).
		Find(&bookmarks).Error

	return bookmarks, total, err
}

func (r *tagRepository) RemoveAllTagsFromBookmark(bookmarkID uuid.UUID) error {
	return r.db.Where("bookmark_id = ?", bookmarkID).Delete(&model.BookmarkTagMapping{}).Error
}

func (r *tagRepository) BulkAddTagToBookmarks(bookmarkIDs []uuid.UUID, tagID uuid.UUID) error {
	for _, bID := range bookmarkIDs {
		mapping := &model.BookmarkTagMapping{
			BookmarkID: bID,
			TagID:      tagID,
		}
		// Ignore duplicate errors
		r.db.Where("bookmark_id = ? AND tag_id = ?", bID, tagID).FirstOrCreate(mapping)
	}
	return nil
}
