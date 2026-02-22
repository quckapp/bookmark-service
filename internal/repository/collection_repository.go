package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type CollectionRepository interface {
	Create(collection *model.BookmarkCollection) error
	GetByID(id uuid.UUID) (*model.BookmarkCollection, error)
	GetByUser(userID uuid.UUID) ([]model.BookmarkCollection, error)
	GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkCollection, error)
	GetPublicByWorkspace(workspaceID uuid.UUID, limit, offset int) ([]model.BookmarkCollection, int64, error)
	Update(collection *model.BookmarkCollection) error
	Delete(id uuid.UUID) error
	AddBookmark(cb *model.CollectionBookmark) error
	RemoveBookmark(collectionID, bookmarkID uuid.UUID) error
	GetBookmarks(collectionID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error)
	IsBookmarkInCollection(collectionID, bookmarkID uuid.UUID) (bool, error)
	CountBookmarks(collectionID uuid.UUID) (int64, error)
}

type collectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) CollectionRepository {
	db.AutoMigrate(&model.BookmarkCollection{}, &model.CollectionBookmark{})
	return &collectionRepository{db: db}
}

func (r *collectionRepository) Create(collection *model.BookmarkCollection) error {
	return r.db.Create(collection).Error
}

func (r *collectionRepository) GetByID(id uuid.UUID) (*model.BookmarkCollection, error) {
	var collection model.BookmarkCollection
	err := r.db.First(&collection, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *collectionRepository) GetByUser(userID uuid.UUID) ([]model.BookmarkCollection, error) {
	var collections []model.BookmarkCollection
	err := r.db.Where("user_id = ?", userID).Order("position ASC, name ASC").Find(&collections).Error
	return collections, err
}

func (r *collectionRepository) GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkCollection, error) {
	var collections []model.BookmarkCollection
	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
		Order("position ASC, name ASC").Find(&collections).Error
	return collections, err
}

func (r *collectionRepository) GetPublicByWorkspace(workspaceID uuid.UUID, limit, offset int) ([]model.BookmarkCollection, int64, error) {
	var collections []model.BookmarkCollection
	var total int64

	r.db.Model(&model.BookmarkCollection{}).Where("workspace_id = ? AND is_public = ?", workspaceID, true).Count(&total)
	err := r.db.Where("workspace_id = ? AND is_public = ?", workspaceID, true).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&collections).Error
	return collections, total, err
}

func (r *collectionRepository) Update(collection *model.BookmarkCollection) error {
	return r.db.Save(collection).Error
}

func (r *collectionRepository) Delete(id uuid.UUID) error {
	r.db.Where("collection_id = ?", id).Delete(&model.CollectionBookmark{})
	return r.db.Delete(&model.BookmarkCollection{}, "id = ?", id).Error
}

func (r *collectionRepository) AddBookmark(cb *model.CollectionBookmark) error {
	return r.db.Create(cb).Error
}

func (r *collectionRepository) RemoveBookmark(collectionID, bookmarkID uuid.UUID) error {
	return r.db.Where("collection_id = ? AND bookmark_id = ?", collectionID, bookmarkID).
		Delete(&model.CollectionBookmark{}).Error
}

func (r *collectionRepository) GetBookmarks(collectionID uuid.UUID, limit, offset int) ([]model.Bookmark, int64, error) {
	var bookmarks []model.Bookmark
	var total int64

	r.db.Model(&model.CollectionBookmark{}).Where("collection_id = ?", collectionID).Count(&total)
	err := r.db.Joins("JOIN collection_bookmarks ON collection_bookmarks.bookmark_id = bookmarks.id").
		Where("collection_bookmarks.collection_id = ?", collectionID).
		Order("collection_bookmarks.position ASC").
		Limit(limit).Offset(offset).
		Find(&bookmarks).Error

	return bookmarks, total, err
}

func (r *collectionRepository) IsBookmarkInCollection(collectionID, bookmarkID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.CollectionBookmark{}).
		Where("collection_id = ? AND bookmark_id = ?", collectionID, bookmarkID).
		Count(&count).Error
	return count > 0, err
}

func (r *collectionRepository) CountBookmarks(collectionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.CollectionBookmark{}).Where("collection_id = ?", collectionID).Count(&count).Error
	return count, err
}
