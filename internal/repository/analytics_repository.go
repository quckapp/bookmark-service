package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetStats(userID uuid.UUID, workspaceID *uuid.UUID) (*model.BookmarkStats, error)
	GetRecentBookmarks(userID uuid.UUID, limit int) ([]model.Bookmark, error)
	GetBookmarkCountByType(userID uuid.UUID) (map[string]int64, error)
	CheckDuplicate(userID uuid.UUID, targetID uuid.UUID, bookmarkType model.BookmarkType) (*model.Bookmark, error)
	SearchBookmarks(userID uuid.UUID, params model.BookmarkSearchParams) ([]model.Bookmark, int64, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetStats(userID uuid.UUID, workspaceID *uuid.UUID) (*model.BookmarkStats, error) {
	stats := &model.BookmarkStats{
		ByType: make(map[string]int64),
	}

	// Total bookmarks
	query := r.db.Model(&model.Bookmark{}).Where("user_id = ?", userID)
	if workspaceID != nil {
		query = query.Where("workspace_id = ?", *workspaceID)
	}
	query.Count(&stats.TotalBookmarks)

	// Total folders
	folderQuery := r.db.Model(&model.BookmarkFolder{}).Where("user_id = ?", userID)
	if workspaceID != nil {
		folderQuery = folderQuery.Where("workspace_id = ?", *workspaceID)
	}
	folderQuery.Count(&stats.TotalFolders)

	// Total tags
	tagQuery := r.db.Model(&model.BookmarkTag{}).Where("user_id = ?", userID)
	if workspaceID != nil {
		tagQuery = tagQuery.Where("workspace_id = ?", *workspaceID)
	}
	tagQuery.Count(&stats.TotalTags)

	// Total collections
	collQuery := r.db.Model(&model.BookmarkCollection{}).Where("user_id = ?", userID)
	if workspaceID != nil {
		collQuery = collQuery.Where("workspace_id = ?", *workspaceID)
	}
	collQuery.Count(&stats.TotalCollections)

	// By type
	type typeCount struct {
		Type  string
		Count int64
	}
	var typeCounts []typeCount
	byTypeQuery := r.db.Model(&model.Bookmark{}).Select("type, count(*) as count").Where("user_id = ?", userID)
	if workspaceID != nil {
		byTypeQuery = byTypeQuery.Where("workspace_id = ?", *workspaceID)
	}
	byTypeQuery.Group("type").Scan(&typeCounts)
	for _, tc := range typeCounts {
		stats.ByType[tc.Type] = tc.Count
	}

	// Recent (last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	recentQuery := r.db.Model(&model.Bookmark{}).Where("user_id = ? AND created_at >= ?", userID, sevenDaysAgo)
	if workspaceID != nil {
		recentQuery = recentQuery.Where("workspace_id = ?", *workspaceID)
	}
	recentQuery.Count(&stats.RecentCount)

	// Favorites
	r.db.Model(&model.BookmarkFavorite{}).Where("user_id = ?", userID).Count(&stats.FavoritesCount)

	return stats, nil
}

func (r *analyticsRepository) GetRecentBookmarks(userID uuid.UUID, limit int) ([]model.Bookmark, error) {
	var bookmarks []model.Bookmark
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&bookmarks).Error
	return bookmarks, err
}

func (r *analyticsRepository) GetBookmarkCountByType(userID uuid.UUID) (map[string]int64, error) {
	result := make(map[string]int64)
	type typeCount struct {
		Type  string
		Count int64
	}
	var counts []typeCount
	err := r.db.Model(&model.Bookmark{}).
		Select("type, count(*) as count").
		Where("user_id = ?", userID).
		Group("type").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}
	for _, c := range counts {
		result[c.Type] = c.Count
	}
	return result, nil
}

func (r *analyticsRepository) CheckDuplicate(userID uuid.UUID, targetID uuid.UUID, bookmarkType model.BookmarkType) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	err := r.db.Where("user_id = ? AND target_id = ? AND type = ?", userID, targetID, bookmarkType).
		First(&bookmark).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &bookmark, nil
}

func (r *analyticsRepository) SearchBookmarks(userID uuid.UUID, params model.BookmarkSearchParams) ([]model.Bookmark, int64, error) {
	var bookmarks []model.Bookmark
	var total int64

	query := r.db.Model(&model.Bookmark{}).Where("user_id = ?", userID)

	if params.Query != "" {
		searchTerm := "%" + params.Query + "%"
		query = query.Where("(title LIKE ? OR description LIKE ?)", searchTerm, searchTerm)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.FolderID != "" {
		folderID, _ := uuid.Parse(params.FolderID)
		query = query.Where("folder_id = ?", folderID)
	}

	query.Count(&total)

	// Sorting
	switch params.Sort {
	case "oldest":
		query = query.Order("created_at ASC")
	case "title":
		query = query.Order("title ASC")
	case "position":
		query = query.Order("position ASC")
	default:
		query = query.Order("created_at DESC")
	}

	offset := params.Page * params.Limit
	err := query.Limit(params.Limit).Offset(offset).Find(&bookmarks).Error
	return bookmarks, total, err
}
