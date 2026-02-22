package service

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/repository"
	"go.uber.org/zap"
)

type TemplateService interface {
	Create(template *model.BookmarkTemplate) error
	GetByUser(userID uuid.UUID) ([]model.BookmarkTemplate, error)
	GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTemplate, error)
	Update(id uuid.UUID, template *model.BookmarkTemplate) (*model.BookmarkTemplate, error)
	Delete(id uuid.UUID) error
	Apply(templateID uuid.UUID, req *model.ApplyTemplateRequest) (*model.Bookmark, error)
}

type templateService struct {
	repo         repository.TemplateRepository
	bookmarkRepo repository.BookmarkRepository
	logger       *zap.Logger
}

func NewTemplateService(repo repository.TemplateRepository, bookmarkRepo repository.BookmarkRepository, logger *zap.Logger) TemplateService {
	return &templateService{repo: repo, bookmarkRepo: bookmarkRepo, logger: logger}
}

func (s *templateService) Create(template *model.BookmarkTemplate) error {
	err := s.repo.Create(template)
	if err != nil {
		return err
	}
	s.logger.Info("Created template", zap.String("id", template.ID.String()))
	return nil
}

func (s *templateService) GetByUser(userID uuid.UUID) ([]model.BookmarkTemplate, error) {
	return s.repo.GetByUser(userID)
}

func (s *templateService) GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTemplate, error) {
	return s.repo.GetByUserAndWorkspace(userID, workspaceID)
}

func (s *templateService) Update(id uuid.UUID, template *model.BookmarkTemplate) (*model.BookmarkTemplate, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if template.Name != "" {
		existing.Name = template.Name
	}
	if template.Description != "" {
		existing.Description = template.Description
	}
	if template.Type != "" {
		existing.Type = template.Type
	}
	if template.FolderID != "" {
		existing.FolderID = template.FolderID
	}
	if template.TagIDs != "" {
		existing.TagIDs = template.TagIDs
	}
	if template.Metadata != "" {
		existing.Metadata = template.Metadata
	}
	err = s.repo.Update(existing)
	return existing, err
}

func (s *templateService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *templateService) Apply(templateID uuid.UUID, req *model.ApplyTemplateRequest) (*model.Bookmark, error) {
	template, err := s.repo.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	targetID, _ := uuid.Parse(req.TargetID)
	var folderID *uuid.UUID
	if template.FolderID != "" {
		fID, _ := uuid.Parse(template.FolderID)
		folderID = &fID
	}

	bookmark := &model.Bookmark{
		UserID:      template.UserID,
		WorkspaceID: template.WorkspaceID,
		FolderID:    folderID,
		Type:        template.Type,
		Title:       req.Title,
		Description: req.Description,
		TargetID:    targetID,
		TargetURL:   req.TargetURL,
		Metadata:    template.Metadata,
	}

	err = s.bookmarkRepo.Create(bookmark)
	if err != nil {
		return nil, err
	}
	s.logger.Info("Applied template", zap.String("templateId", templateID.String()))
	return bookmark, nil
}
