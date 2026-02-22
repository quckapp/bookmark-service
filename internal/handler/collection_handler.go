package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/service"
)

type CollectionHandler struct {
	service service.CollectionService
}

func NewCollectionHandler(service service.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}

func (h *CollectionHandler) Create(c *gin.Context) {
	var req model.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(req.UserID)
	workspaceID, _ := uuid.Parse(req.WorkspaceID)

	collection := &model.BookmarkCollection{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsPublic:    req.IsPublic,
	}

	if err := h.service.Create(collection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": collection})
}

func (h *CollectionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collection})
}

func (h *CollectionHandler) GetByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	collections, err := h.service.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collections})
}

func (h *CollectionHandler) GetByUserAndWorkspace(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("userId"))
	workspaceID, _ := uuid.Parse(c.Param("workspaceId"))

	collections, err := h.service.GetByUserAndWorkspace(userID, workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collections})
}

func (h *CollectionHandler) GetPublic(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	collections, total, err := h.service.GetPublicByWorkspace(workspaceID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collections, "total": total, "page": page, "limit": limit})
}

func (h *CollectionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req model.UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection, err := h.service.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collection})
}

func (h *CollectionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted"})
}

func (h *CollectionHandler) AddBookmarks(c *gin.Context) {
	collectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req model.AddToCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bookmarkIDs []uuid.UUID
	for _, id := range req.BookmarkIDs {
		bID, _ := uuid.Parse(id)
		bookmarkIDs = append(bookmarkIDs, bID)
	}

	if err := h.service.AddBookmarks(collectionID, bookmarkIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmarks added to collection"})
}

func (h *CollectionHandler) RemoveBookmark(c *gin.Context) {
	collectionID, _ := uuid.Parse(c.Param("id"))
	bookmarkID, _ := uuid.Parse(c.Param("bookmarkId"))

	if err := h.service.RemoveBookmark(collectionID, bookmarkID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark removed from collection"})
}

func (h *CollectionHandler) GetBookmarks(c *gin.Context) {
	collectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	bookmarks, total, err := h.service.GetBookmarks(collectionID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookmarks, "total": total, "page": page, "limit": limit})
}
