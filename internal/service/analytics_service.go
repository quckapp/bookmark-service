package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type BookmarkAnalyticsService interface {
	GetStats(userID uuid.UUID, workspaceID *uuid.UUID) (*model.BookmarkStats, error)
	GetRecentBookmarks(userID uuid.UUID, limit int) ([]model.Bookmark, error)
	CheckDuplicate(userID uuid.UUID, targetID uuid.UUID, bookmarkType model.BookmarkType) (*model.DuplicateCheckResult, error)
	SearchBookmarks(userID uuid.UUID, params model.BookmarkSearchParams) ([]model.Bookmark, int64, error)
	ExportBookmarks(userID uuid.UUID, workspaceID *uuid.UUID) (*model.ExportData, error)
	ImportBookmarks(userID, workspaceID uuid.UUID, req model.ImportRequest) (*model.ImportResult, error)
	GetActivity(userID uuid.UUID, page, limit int) ([]model.BookmarkActivity, int64, error)
	GetBookmarkActivity(bookmarkID uuid.UUID, page, limit int) ([]model.BookmarkActivity, int64, error)
	LogActivity(bookmarkID, userID uuid.UUID, action, details string) error
}

type bookmarkAnalyticsService struct {
	analyticsRepo repository.AnalyticsRepository
	bookmarkRepo  repository.BookmarkRepository
	folderRepo    repository.FolderRepository
	tagRepo       repository.TagRepository
	collectionRepo repository.CollectionRepository
	activityRepo  repository.ActivityRepository
	logger        *zap.Logger
}

func NewBookmarkAnalyticsService(
	analyticsRepo repository.AnalyticsRepository,
	bookmarkRepo repository.BookmarkRepository,
	folderRepo repository.FolderRepository,
	tagRepo repository.TagRepository,
	collectionRepo repository.CollectionRepository,
	activityRepo repository.ActivityRepository,
	logger *zap.Logger,
) BookmarkAnalyticsService {
	return &bookmarkAnalyticsService{
		analyticsRepo:  analyticsRepo,
		bookmarkRepo:   bookmarkRepo,
		folderRepo:     folderRepo,
		tagRepo:        tagRepo,
		collectionRepo: collectionRepo,
		activityRepo:   activityRepo,
		logger:         logger,
	}
}

func (s *bookmarkAnalyticsService) GetStats(userID uuid.UUID, workspaceID *uuid.UUID) (*model.BookmarkStats, error) {
	return s.analyticsRepo.GetStats(userID, workspaceID)
}

func (s *bookmarkAnalyticsService) GetRecentBookmarks(userID uuid.UUID, limit int) ([]model.Bookmark, error) {
	return s.analyticsRepo.GetRecentBookmarks(userID, limit)
}

func (s *bookmarkAnalyticsService) CheckDuplicate(userID uuid.UUID, targetID uuid.UUID, bookmarkType model.BookmarkType) (*model.DuplicateCheckResult, error) {
	existing, err := s.analyticsRepo.CheckDuplicate(userID, targetID, bookmarkType)
	if err != nil {
		return nil, err
	}
	result := &model.DuplicateCheckResult{
		IsDuplicate: existing != nil,
		Existing:    existing,
	}
	return result, nil
}

func (s *bookmarkAnalyticsService) SearchBookmarks(userID uuid.UUID, params model.BookmarkSearchParams) ([]model.Bookmark, int64, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	return s.analyticsRepo.SearchBookmarks(userID, params)
}

func (s *bookmarkAnalyticsService) ExportBookmarks(userID uuid.UUID, workspaceID *uuid.UUID) (*model.ExportData, error) {
	export := &model.ExportData{
		ExportedAt: time.Now(),
	}

	// Get bookmarks
	bookmarks, _, err := s.bookmarkRepo.GetByUser(userID, 10000, 0)
	if err != nil {
		return nil, err
	}
	export.Bookmarks = bookmarks

	// Get folders
	folders, err := s.folderRepo.GetByUser(userID)
	if err != nil {
		return nil, err
	}
	export.Folders = folders

	// Get tags
	tags, err := s.tagRepo.GetByUser(userID)
	if err != nil {
		return nil, err
	}
	export.Tags = tags

	// Get collections
	collections, err := s.collectionRepo.GetByUser(userID)
	if err != nil {
		return nil, err
	}
	export.Collections = collections

	return export, nil
}

func (s *bookmarkAnalyticsService) ImportBookmarks(userID, workspaceID uuid.UUID, req model.ImportRequest) (*model.ImportResult, error) {
	result := &model.ImportResult{}

	var folderID *uuid.UUID
	if req.FolderID != "" {
		fID, _ := uuid.Parse(req.FolderID)
		folderID = &fID
	}

	for _, item := range req.Bookmarks {
		bookmark := &model.Bookmark{
			UserID:      userID,
			WorkspaceID: workspaceID,
			FolderID:    folderID,
			Type:        model.BookmarkType(item.Type),
			Title:       item.Title,
			Description: item.Description,
			TargetID:    uuid.New(), // Generate target ID for imports
			TargetURL:   item.TargetURL,
			Metadata:    item.Metadata,
		}

		if err := s.bookmarkRepo.Create(bookmark); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Imported++
		}
	}

	return result, nil
}

func (s *bookmarkAnalyticsService) GetActivity(userID uuid.UUID, page, limit int) ([]model.BookmarkActivity, int64, error) {
	offset := page * limit
	return s.activityRepo.GetByUser(userID, limit, offset)
}

func (s *bookmarkAnalyticsService) GetBookmarkActivity(bookmarkID uuid.UUID, page, limit int) ([]model.BookmarkActivity, int64, error) {
	offset := page * limit
	return s.activityRepo.GetByBookmark(bookmarkID, limit, offset)
}

func (s *bookmarkAnalyticsService) LogActivity(bookmarkID, userID uuid.UUID, action, details string) error {
	activity := &model.BookmarkActivity{
		BookmarkID: bookmarkID,
		UserID:     userID,
		Action:     action,
		Details:    details,
	}
	return s.activityRepo.Create(activity)
}
