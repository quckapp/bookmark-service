package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/quckapp/bookmark-service/internal/model"
	"github.com/quckapp/bookmark-service/internal/service"
)

type ReminderHandler struct {
	service service.BookmarkReminderService
}

func NewReminderHandler(service service.BookmarkReminderService) *ReminderHandler {
	return &ReminderHandler{service: service}
}

func (h *ReminderHandler) Create(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	var req model.CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	remindAt, err := time.Parse(time.RFC3339, req.RemindAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid remindAt format, use RFC3339"})
		return
	}

	userID, _ := uuid.Parse(c.Param("userId"))

	reminder := &model.BookmarkReminder{
		BookmarkID: bookmarkID,
		UserID:     userID,
		RemindAt:   remindAt,
		Message:    req.Message,
	}

	if err := h.service.Create(reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": reminder})
}

func (h *ReminderHandler) GetByBookmark(c *gin.Context) {
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookmark ID"})
		return
	}

	reminders, err := h.service.GetByBookmark(bookmarkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminders})
}

func (h *ReminderHandler) GetByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	reminders, total, err := h.service.GetByUser(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminders, "total": total, "page": page, "limit": limit})
}

func (h *ReminderHandler) GetPending(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	reminders, err := h.service.GetPending(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminders})
}

func (h *ReminderHandler) Cancel(c *gin.Context) {
	reminderID, err := uuid.Parse(c.Param("reminderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}

	if err := h.service.Cancel(reminderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reminder cancelled"})
}

func (h *ReminderHandler) Delete(c *gin.Context) {
	reminderID, err := uuid.Parse(c.Param("reminderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}

	if err := h.service.Delete(reminderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted"})
}
