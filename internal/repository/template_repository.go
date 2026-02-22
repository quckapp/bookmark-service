package repository

import (
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type TemplateRepository interface {
	Create(template *model.BookmarkTemplate) error
	GetByID(id uuid.UUID) (*model.BookmarkTemplate, error)
	GetByUser(userID uuid.UUID) ([]model.BookmarkTemplate, error)
	GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTemplate, error)
	Update(template *model.BookmarkTemplate) error
	Delete(id uuid.UUID) error
}

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	db.AutoMigrate(&model.BookmarkTemplate{})
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(template *model.BookmarkTemplate) error {
	return r.db.Create(template).Error
}

func (r *templateRepository) GetByID(id uuid.UUID) (*model.BookmarkTemplate, error) {
	var template model.BookmarkTemplate
	err := r.db.First(&template, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) GetByUser(userID uuid.UUID) ([]model.BookmarkTemplate, error) {
	var templates []model.BookmarkTemplate
	err := r.db.Where("user_id = ?", userID).Order("name ASC").Find(&templates).Error
	return templates, err
}

func (r *templateRepository) GetByUserAndWorkspace(userID, workspaceID uuid.UUID) ([]model.BookmarkTemplate, error) {
	var templates []model.BookmarkTemplate
	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
		Order("name ASC").Find(&templates).Error
	return templates, err
}

func (r *templateRepository) Update(template *model.BookmarkTemplate) error {
	return r.db.Save(template).Error
}

func (r *templateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.BookmarkTemplate{}, "id = ?", id).Error
}
