package service

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type VersionService interface {
	GetByBookmark(bookmarkID uuid.UUID, page, limit int) ([]model.BookmarkVersion, int64, error)
	GetByID(id uuid.UUID) (*model.BookmarkVersion, error)
	Restore(bookmarkID, versionID uuid.UUID) (*model.Bookmark, error)
	CreateVersion(bookmark *model.Bookmark, userID uuid.UUID, changeNote string) error
}

type versionService struct {
	repo         repository.VersionRepository
	bookmarkRepo repository.BookmarkRepository
	logger       *zap.Logger
}

func NewVersionService(repo repository.VersionRepository, bookmarkRepo repository.BookmarkRepository, logger *zap.Logger) VersionService {
	return &versionService{repo: repo, bookmarkRepo: bookmarkRepo, logger: logger}
}

func (s *versionService) GetByBookmark(bookmarkID uuid.UUID, page, limit int) ([]model.BookmarkVersion, int64, error) {
	offset := page * limit
	return s.repo.GetByBookmark(bookmarkID, limit, offset)
}

func (s *versionService) GetByID(id uuid.UUID) (*model.BookmarkVersion, error) {
	return s.repo.GetByID(id)
}

func (s *versionService) Restore(bookmarkID, versionID uuid.UUID) (*model.Bookmark, error) {
	version, err := s.repo.GetByID(versionID)
	if err != nil {
		return nil, err
	}

	bookmark, err := s.bookmarkRepo.GetByID(bookmarkID)
	if err != nil {
		return nil, err
	}

	// Save current state as new version before restoring
	s.CreateVersion(bookmark, bookmark.UserID, "Before restore")

	// Restore from version
	bookmark.Title = version.Title
	bookmark.Description = version.Description
	bookmark.TargetURL = version.TargetURL
	bookmark.Metadata = version.Metadata

	err = s.bookmarkRepo.Update(bookmark)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Restored bookmark from version",
		zap.String("bookmarkId", bookmarkID.String()),
		zap.String("versionId", versionID.String()))
	return bookmark, nil
}

func (s *versionService) CreateVersion(bookmark *model.Bookmark, userID uuid.UUID, changeNote string) error {
	latestVer, _ := s.repo.GetLatestVersion(bookmark.ID)

	version := &model.BookmarkVersion{
		BookmarkID:  bookmark.ID,
		UserID:      userID,
		Title:       bookmark.Title,
		Description: bookmark.Description,
		TargetURL:   bookmark.TargetURL,
		Metadata:    bookmark.Metadata,
		Version:     latestVer + 1,
		ChangeNote:  changeNote,
	}
	return s.repo.Create(version)
}
