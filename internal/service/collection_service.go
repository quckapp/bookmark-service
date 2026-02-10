package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type CollectionService interface {
	Create(collection *model.BookmarkCollection) error
	GetByID(id uuid.UUID) (*model.BookmarkCollection, error)
	GetByUser(userID uuid.UUID) ([]model.BookmarkCollection, error)
	GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkCollection, error)
	GetPublicByWorkspace(workspaceID uuid.UUID, page, limit int) ([]model.BookmarkCollection, int64, error)
	Update(id uuid.UUID, req *model.UpdateCollectionRequest) (*model.BookmarkCollection, error)
	Delete(id uuid.UUID) error
	AddBookmarks(collectionID uuid.UUID, bookmarkIDs []uuid.UUID) error
	RemoveBookmark(collectionID, bookmarkID uuid.UUID) error
	GetBookmarks(collectionID uuid.UUID, page, limit int) ([]model.Bookmark, int64, error)
	CountBookmarks(collectionID uuid.UUID) (int64, error)
}

type collectionService struct {
	repo   repository.CollectionRepository
	logger *zap.Logger
}

func NewCollectionService(repo repository.CollectionRepository, logger *zap.Logger) CollectionService {
	return &collectionService{repo: repo, logger: logger}
}

func (s *collectionService) Create(collection *model.BookmarkCollection) error {
	err := s.repo.Create(collection)
	if err != nil {
		return err
	}
	s.logger.Info("Created collection", zap.String("id", collection.ID.String()))
	return nil
}

func (s *collectionService) GetByID(id uuid.UUID) (*model.BookmarkCollection, error) {
	return s.repo.GetByID(id)
}

func (s *collectionService) GetByUser(userID uuid.UUID) ([]model.BookmarkCollection, error) {
	return s.repo.GetByUser(userID)
}

func (s *collectionService) GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkCollection, error) {
	return s.repo.GetByUserAndWorkspace(userID, workspaceID)
}

func (s *collectionService) GetPublicByWorkspace(workspaceID uuid.UUID, page, limit int) ([]model.BookmarkCollection, int64, error) {
	offset := page * limit
	return s.repo.GetPublicByWorkspace(workspaceID, limit, offset)
}

func (s *collectionService) Update(id uuid.UUID, req *model.UpdateCollectionRequest) (*model.BookmarkCollection, error) {
	collection, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("collection not found")
	}

	if req.Name != "" {
		collection.Name = req.Name
	}
	if req.Description != "" {
		collection.Description = req.Description
	}
	if req.Color != "" {
		collection.Color = req.Color
	}
	if req.Icon != "" {
		collection.Icon = req.Icon
	}
	if req.IsPublic != nil {
		collection.IsPublic = *req.IsPublic
	}
	if req.Position != nil {
		collection.Position = *req.Position
	}

	err = s.repo.Update(collection)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (s *collectionService) Delete(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	s.logger.Info("Deleted collection", zap.String("id", id.String()))
	return nil
}

func (s *collectionService) AddBookmarks(collectionID uuid.UUID, bookmarkIDs []uuid.UUID) error {
	for _, bID := range bookmarkIDs {
		exists, _ := s.repo.IsBookmarkInCollection(collectionID, bID)
		if exists {
			continue
		}
		cb := &model.CollectionBookmark{
			CollectionID: collectionID,
			BookmarkID:   bID,
		}
		if err := s.repo.AddBookmark(cb); err != nil {
			s.logger.Warn("Failed to add bookmark to collection",
				zap.String("collectionID", collectionID.String()),
				zap.String("bookmarkID", bID.String()),
				zap.Error(err))
		}
	}
	return nil
}

func (s *collectionService) RemoveBookmark(collectionID, bookmarkID uuid.UUID) error {
	return s.repo.RemoveBookmark(collectionID, bookmarkID)
}

func (s *collectionService) GetBookmarks(collectionID uuid.UUID, page, limit int) ([]model.Bookmark, int64, error) {
	offset := page * limit
	return s.repo.GetBookmarks(collectionID, limit, offset)
}

func (s *collectionService) CountBookmarks(collectionID uuid.UUID) (int64, error) {
	return s.repo.CountBookmarks(collectionID)
}
