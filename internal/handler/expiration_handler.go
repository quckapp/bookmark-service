package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/service"
)

type ExpirationHandler struct {
	service service.ExpirationService
}

func NewExpirationHandler(service service.ExpirationService) *ExpirationHandler {
	return &ExpirationHandler{service: service}
}

func (h *ExpirationHandler) Set(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req model.SetExpirationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiresAt format, use RFC3339"})
		return
	}

	action := req.Action
	if action == "" {
		action = "archive"
	}

	expiration := &model.BookmarkExpiration{
		BookmarkID: bookmarkID,
		UserID:     userID,
		ExpiresAt:  expiresAt,
		Action:     action,
	}

	if err := h.service.Set(expiration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": expiration})
}

func (h *ExpirationHandler) Get(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	expiration, err := h.service.GetByBookmarkID(bookmarkID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expiration not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expiration})
}

func (h *ExpirationHandler) Remove(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	if err := h.service.Remove(bookmarkID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expiration removed"})
}

func (h *ExpirationHandler) GetExpiring(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	expirations, err := h.service.GetExpiring(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expirations})
}
