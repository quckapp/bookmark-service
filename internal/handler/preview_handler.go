package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/service"
)

type PreviewHandler struct {
	service service.PreviewService
}

func NewPreviewHandler(service service.PreviewService) *PreviewHandler {
	return &PreviewHandler{service: service}
}

func (h *PreviewHandler) Generate(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	preview, err := h.service.Generate(bookmarkID, req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": preview})
}

func (h *PreviewHandler) Get(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	preview, err := h.service.GetByBookmarkID(bookmarkID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preview not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": preview})
}

func (h *PreviewHandler) Create(c *gin.Context) {
	var req model.CreatePreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookmarkID, _ := uuid.Parse(req.BookmarkID)
	preview := &model.LinkPreview{
		BookmarkID: bookmarkID,
		URL:        req.URL,
		FetchedAt:  time.Now(),
	}

	if err := h.service.Create(preview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": preview})
}

func (h *PreviewHandler) GetByURL(c *gin.Context) {
	url := c.Param("url")
	preview, err := h.service.GetByURL(url)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preview not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": preview})
}
