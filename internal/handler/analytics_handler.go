package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/service"
)

type AnalyticsHandler struct {
	service service.BookmarkAnalyticsService
}

func NewAnalyticsHandler(service service.BookmarkAnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetStats(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var workspaceID *uuid.UUID
	if wsID := c.Query("workspaceId"); wsID != "" {
		parsed, _ := uuid.Parse(wsID)
		workspaceID = &parsed
	}

	stats, err := h.service.GetStats(userID, workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *AnalyticsHandler) GetRecent(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	bookmarks, err := h.service.GetRecentBookmarks(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookmarks})
}

func (h *AnalyticsHandler) CheckDuplicate(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("userId"))
	targetID, _ := uuid.Parse(c.Query("targetId"))
	bookmarkType := model.BookmarkType(c.Query("type"))

	result, err := h.service.CheckDuplicate(userID, targetID, bookmarkType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *AnalyticsHandler) Search(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var params model.BookmarkSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookmarks, total, err := h.service.SearchBookmarks(userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookmarks, "total": total, "page": params.Page, "limit": params.Limit})
}

func (h *AnalyticsHandler) Export(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var workspaceID *uuid.UUID
	if wsID := c.Query("workspaceId"); wsID != "" {
		parsed, _ := uuid.Parse(wsID)
		workspaceID = &parsed
	}

	data, err := h.service.ExportBookmarks(userID, workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *AnalyticsHandler) Import(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("userId"))
	workspaceID, _ := uuid.Parse(c.Param("workspaceId"))

	var req model.ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ImportBookmarks(userID, workspaceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *AnalyticsHandler) GetActivity(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	activities, total, err := h.service.GetActivity(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": activities, "total": total, "page": page, "limit": limit})
}

func (h *AnalyticsHandler) GetBookmarkActivity(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	activities, total, err := h.service.GetBookmarkActivity(bookmarkID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": activities, "total": total, "page": page, "limit": limit})
}
