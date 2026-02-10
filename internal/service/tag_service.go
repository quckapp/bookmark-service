package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type TagService interface {
	CreateTag(tag *model.BookmarkTag) error
	GetTag(id uuid.UUID) (*model.BookmarkTag, error)
	GetUserTags(userID uuid.UUID) ([]model.BookmarkTag, error)
	GetUserTagsInWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTag, error)
	UpdateTag(id uuid.UUID, req *model.UpdateTagRequest) (*model.BookmarkTag, error)
	DeleteTag(id uuid.UUID) error
	TagBookmark(bookmarkID uuid.UUID, tagIDs []uuid.UUID) error
	UntagBookmark(bookmarkID, tagID uuid.UUID) error
	GetBookmarkTags(bookmarkID uuid.UUID) ([]model.BookmarkTag, error)
	GetBookmarksByTag(tagID uuid.UUID, page, limit int) ([]model.Bookmark, int64, error)
	ReplaceBookmarkTags(bookmarkID uuid.UUID, tagIDs []uuid.UUID) error
	BulkTagBookmarks(bookmarkIDs []uuid.UUID, tagID uuid.UUID) error
}

type tagService struct {
	repo   repository.TagRepository
	logger *zap.Logger
}

func NewTagService(repo repository.TagRepository, logger *zap.Logger) TagService {
	return &tagService{repo: repo, logger: logger}
}

func (s *tagService) CreateTag(tag *model.BookmarkTag) error {
	err := s.repo.Create(tag)
	if err != nil {
		return err
	}
	s.logger.Info("Created tag", zap.String("id", tag.ID.String()))
	return nil
}

func (s *tagService) GetTag(id uuid.UUID) (*model.BookmarkTag, error) {
	return s.repo.GetByID(id)
}

func (s *tagService) GetUserTags(userID uuid.UUID) ([]model.BookmarkTag, error) {
	return s.repo.GetByUser(userID)
}

func (s *tagService) GetUserTagsInWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTag, error) {
	return s.repo.GetByUserAndWorkspace(userID, workspaceID)
}

func (s *tagService) UpdateTag(id uuid.UUID, req *model.UpdateTagRequest) (*model.BookmarkTag, error) {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("tag not found")
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}

	err = s.repo.Update(tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *tagService) DeleteTag(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	s.logger.Info("Deleted tag", zap.String("id", id.String()))
	return nil
}

func (s *tagService) TagBookmark(bookmarkID uuid.UUID, tagIDs []uuid.UUID) error {
	for _, tagID := range tagIDs {
		mapping := &model.BookmarkTagMapping{
			BookmarkID: bookmarkID,
			TagID:      tagID,
		}
		if err := s.repo.AddTagToBookmark(mapping); err != nil {
			s.logger.Warn("Failed to add tag to bookmark",
				zap.String("bookmarkID", bookmarkID.String()),
				zap.String("tagID", tagID.String()),
				zap.Error(err))
		}
	}
	return nil
}

func (s *tagService) UntagBookmark(bookmarkID, tagID uuid.UUID) error {
	return s.repo.RemoveTagFromBookmark(bookmarkID, tagID)
}

func (s *tagService) GetBookmarkTags(bookmarkID uuid.UUID) ([]model.BookmarkTag, error) {
	return s.repo.GetTagsByBookmark(bookmarkID)
}

func (s *tagService) GetBookmarksByTag(tagID uuid.UUID, page, limit int) ([]model.Bookmark, int64, error) {
	offset := page * limit
	return s.repo.GetBookmarksByTag(tagID, limit, offset)
}

func (s *tagService) ReplaceBookmarkTags(bookmarkID uuid.UUID, tagIDs []uuid.UUID) error {
	if err := s.repo.RemoveAllTagsFromBookmark(bookmarkID); err != nil {
		return err
	}
	return s.TagBookmark(bookmarkID, tagIDs)
}

func (s *tagService) BulkTagBookmarks(bookmarkIDs []uuid.UUID, tagID uuid.UUID) error {
	return s.repo.BulkAddTagToBookmarks(bookmarkIDs, tagID)
}
