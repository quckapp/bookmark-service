package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bookmark struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	WorkspaceID uuid.UUID      `gorm:"type:char(36);not null;index" json:"workspaceId"`
	FolderID    *uuid.UUID     `gorm:"type:char(36);index" json:"folderId,omitempty"`
	Type        BookmarkType   `gorm:"type:varchar(20);not null" json:"type"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	TargetID    uuid.UUID      `gorm:"type:char(36);not null" json:"targetId"`
	TargetURL   string         `gorm:"type:varchar(500)" json:"targetUrl,omitempty"`
	Metadata    string         `gorm:"type:json" json:"metadata,omitempty"`
	Position    int            `gorm:"default:0" json:"position"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type BookmarkType string

const (
	BookmarkTypeMessage  BookmarkType = "message"
	BookmarkTypeChannel  BookmarkType = "channel"
	BookmarkTypeFile     BookmarkType = "file"
	BookmarkTypeThread   BookmarkType = "thread"
	BookmarkTypeExternal BookmarkType = "external"
)

func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BookmarkFolder struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	WorkspaceID uuid.UUID      `gorm:"type:char(36);not null;index" json:"workspaceId"`
	ParentID    *uuid.UUID     `gorm:"type:char(36);index" json:"parentId,omitempty"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Color       string         `gorm:"type:varchar(20)" json:"color,omitempty"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon,omitempty"`
	Position    int            `gorm:"default:0" json:"position"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (f *BookmarkFolder) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type BookmarkTag struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	WorkspaceID uuid.UUID      `gorm:"type:char(36);not null;index" json:"workspaceId"`
	Name        string         `gorm:"type:varchar(50);not null" json:"name"`
	Color       string         `gorm:"type:varchar(20)" json:"color,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (bt *BookmarkTag) BeforeCreate(tx *gorm.DB) error {
	if bt.ID == uuid.Nil {
		bt.ID = uuid.New()
	}
	return nil
}

type BookmarkTagMapping struct {
	ID         uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID uuid.UUID `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	TagID      uuid.UUID `gorm:"type:char(36);not null;index" json:"tagId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (m *BookmarkTagMapping) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type BookmarkCollection struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	WorkspaceID uuid.UUID      `gorm:"type:char(36);not null;index" json:"workspaceId"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Color       string         `gorm:"type:varchar(20)" json:"color,omitempty"`
	Icon        string         `gorm:"type:varchar(50)" json:"icon,omitempty"`
	IsPublic    bool           `gorm:"default:false" json:"isPublic"`
	Position    int            `gorm:"default:0" json:"position"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (bc *BookmarkCollection) BeforeCreate(tx *gorm.DB) error {
	if bc.ID == uuid.Nil {
		bc.ID = uuid.New()
	}
	return nil
}

type CollectionBookmark struct {
	ID           uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	CollectionID uuid.UUID `gorm:"type:char(36);not null;index" json:"collectionId"`
	BookmarkID   uuid.UUID `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	Position     int       `gorm:"default:0" json:"position"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (cb *CollectionBookmark) BeforeCreate(tx *gorm.DB) error {
	if cb.ID == uuid.Nil {
		cb.ID = uuid.New()
	}
	return nil
}

type SharedBookmark struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID  uuid.UUID      `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	SharedBy    uuid.UUID      `gorm:"type:char(36);not null;index" json:"sharedBy"`
	SharedWith  uuid.UUID      `gorm:"type:char(36);not null;index" json:"sharedWith"`
	WorkspaceID uuid.UUID      `gorm:"type:char(36);not null" json:"workspaceId"`
	Message     string         `gorm:"type:text" json:"message,omitempty"`
	IsAccepted  bool           `gorm:"default:false" json:"isAccepted"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (sb *SharedBookmark) BeforeCreate(tx *gorm.DB) error {
	if sb.ID == uuid.Nil {
		sb.ID = uuid.New()
	}
	return nil
}

type BookmarkNote struct {
	ID         uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID uuid.UUID      `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	UserID     uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Color      string         `gorm:"type:varchar(20)" json:"color,omitempty"`
	IsPinned   bool           `gorm:"default:false" json:"isPinned"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (bn *BookmarkNote) BeforeCreate(tx *gorm.DB) error {
	if bn.ID == uuid.Nil {
		bn.ID = uuid.New()
	}
	return nil
}

type BookmarkReminder struct {
	ID         uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID uuid.UUID      `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	UserID     uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	RemindAt   time.Time      `gorm:"not null" json:"remindAt"`
	Message    string         `gorm:"type:text" json:"message,omitempty"`
	Status     string         `gorm:"type:varchar(20);default:pending" json:"status"`
	FiredAt    *time.Time     `json:"firedAt,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (br *BookmarkReminder) BeforeCreate(tx *gorm.DB) error {
	if br.ID == uuid.Nil {
		br.ID = uuid.New()
	}
	return nil
}

type BookmarkFavorite struct {
	ID         uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	UserID     uuid.UUID `gorm:"type:char(36);not null;index" json:"userId"`
	BookmarkID uuid.UUID `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (bf *BookmarkFavorite) BeforeCreate(tx *gorm.DB) error {
	if bf.ID == uuid.Nil {
		bf.ID = uuid.New()
	}
	return nil
}

type BookmarkActivity struct {
	ID         uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID uuid.UUID `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	UserID     uuid.UUID `gorm:"type:char(36);not null;index" json:"userId"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`
	Details    string    `gorm:"type:text" json:"details,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (ba *BookmarkActivity) BeforeCreate(tx *gorm.DB) error {
	if ba.ID == uuid.Nil {
		ba.ID = uuid.New()
	}
	return nil
}

type LinkPreview struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID  uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex" json:"bookmarkId"`
	URL         string         `gorm:"type:varchar(500);not null" json:"url"`
	Title       string         `gorm:"type:varchar(255)" json:"title,omitempty"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"imageUrl,omitempty"`
	FaviconURL  string         `gorm:"type:varchar(500)" json:"faviconUrl,omitempty"`
	SiteName    string         `gorm:"type:varchar(100)" json:"siteName,omitempty"`
	ContentType string         `gorm:"type:varchar(50)" json:"contentType,omitempty"`
	FetchedAt   time.Time      `json:"fetchedAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (lp *LinkPreview) BeforeCreate(tx *gorm.DB) error {
	if lp.ID == uuid.Nil {
		lp.ID = uuid.New()
	}
	return nil
}

type ReadLaterStatus string

const (
	ReadLaterStatusUnread    ReadLaterStatus = "unread"
	ReadLaterStatusReading   ReadLaterStatus = "reading"
	ReadLaterStatusCompleted ReadLaterStatus = "completed"
	ReadLaterStatusArchived  ReadLaterStatus = "archived"
)

type ReadLaterItem struct {
	ID         uuid.UUID       `gorm:"type:char(36);primary_key" json:"id"`
	UserID     uuid.UUID       `gorm:"type:char(36);not null;index" json:"userId"`
	BookmarkID uuid.UUID       `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	Status     ReadLaterStatus `gorm:"type:varchar(20);default:unread" json:"status"`
	Priority   int             `gorm:"default:0" json:"priority"`
	ReadAt     *time.Time      `json:"readAt,omitempty"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (rl *ReadLaterItem) BeforeCreate(tx *gorm.DB) error {
	if rl.ID == uuid.Nil {
		rl.ID = uuid.New()
	}
	return nil
}

type BookmarkComment struct {
	ID         uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID uuid.UUID      `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	UserID     uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	ParentID   *uuid.UUID     `gorm:"type:char(36);index" json:"parentId,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (bc *BookmarkComment) BeforeCreate(tx *gorm.DB) error {
	if bc.ID == uuid.Nil {
		bc.ID = uuid.New()
	}
	return nil
}

type BookmarkVersion struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID  uuid.UUID `gorm:"type:char(36);not null;index" json:"bookmarkId"`
	UserID      uuid.UUID `gorm:"type:char(36);not null" json:"userId"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	TargetURL   string    `gorm:"type:varchar(500)" json:"targetUrl,omitempty"`
	Metadata    string    `gorm:"type:json" json:"metadata,omitempty"`
	Version     int       `gorm:"not null" json:"version"`
	ChangeNote  string    `gorm:"type:text" json:"changeNote,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (bv *BookmarkVersion) BeforeCreate(tx *gorm.DB) error {
	if bv.ID == uuid.Nil {
		bv.ID = uuid.New()
	}
	return nil
}

type BookmarkExpiration struct {
	ID         uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	BookmarkID uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex" json:"bookmarkId"`
	UserID     uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	ExpiresAt  time.Time      `gorm:"not null" json:"expiresAt"`
	Action     string         `gorm:"type:varchar(20);default:archive" json:"action"`
	IsExpired  bool           `gorm:"default:false" json:"isExpired"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (be *BookmarkExpiration) BeforeCreate(tx *gorm.DB) error {
	if be.ID == uuid.Nil {
		be.ID = uuid.New()
	}
	return nil
}

type BookmarkTemplate struct {
	ID          uuid.UUID      `gorm:"type:char(36);primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"userId"`
	WorkspaceID uuid.UUID      `gorm:"type:char(36);not null;index" json:"workspaceId"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Type        BookmarkType   `gorm:"type:varchar(20)" json:"type,omitempty"`
	FolderID    string         `gorm:"type:char(36)" json:"folderId,omitempty"`
	TagIDs      string         `gorm:"type:json" json:"tagIds,omitempty"`
	Metadata    string         `gorm:"type:json" json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (bt *BookmarkTemplate) BeforeCreate(tx *gorm.DB) error {
	if bt.ID == uuid.Nil {
		bt.ID = uuid.New()
	}
	return nil
}

// Analytics & Stats DTOs

type BookmarkStats struct {
	TotalBookmarks   int64            `json:"totalBookmarks"`
	TotalFolders     int64            `json:"totalFolders"`
	TotalTags        int64            `json:"totalTags"`
	TotalCollections int64            `json:"totalCollections"`
	ByType           map[string]int64 `json:"byType"`
	RecentCount      int64            `json:"recentCount"`
	FavoritesCount   int64            `json:"favoritesCount"`
}

type DuplicateCheckResult struct {
	IsDuplicate bool      `json:"isDuplicate"`
	Existing    *Bookmark `json:"existing,omitempty"`
}

type BookmarkSearchParams struct {
	Query    string `form:"query" json:"query"`
	Type     string `form:"type" json:"type,omitempty"`
	FolderID string `form:"folderId" json:"folderId,omitempty"`
	Sort     string `form:"sort" json:"sort,omitempty"`
	Page     int    `form:"page" json:"page"`
	Limit    int    `form:"limit" json:"limit"`
}

type ExportData struct {
	Bookmarks   []Bookmark           `json:"bookmarks"`
	Folders     []BookmarkFolder     `json:"folders"`
	Tags        []BookmarkTag        `json:"tags"`
	Collections []BookmarkCollection `json:"collections"`
	ExportedAt  time.Time            `json:"exportedAt"`
}

type ImportRequest struct {
	Bookmarks []ImportItem `json:"bookmarks" binding:"required"`
	FolderID  string       `json:"folderId,omitempty"`
}

type ImportItem struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	TargetURL   string `json:"targetUrl,omitempty"`
	Metadata    string `json:"metadata,omitempty"`
}

type ImportResult struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

type ReadLaterStats struct {
	Total     int64 `json:"total"`
	Unread    int64 `json:"unread"`
	Reading   int64 `json:"reading"`
	Completed int64 `json:"completed"`
	Archived  int64 `json:"archived"`
}

// Request DTOs

type CreateTagRequest struct {
	UserID      string `json:"userId" binding:"required"`
	WorkspaceID string `json:"workspaceId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Color       string `json:"color,omitempty"`
}

type UpdateTagRequest struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

type TagBookmarkRequest struct {
	TagIDs []string `json:"tagIds" binding:"required"`
}

type CreateCollectionRequest struct {
	UserID      string `json:"userId" binding:"required"`
	WorkspaceID string `json:"workspaceId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	IsPublic    bool   `json:"isPublic"`
}

type UpdateCollectionRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	IsPublic    *bool  `json:"isPublic,omitempty"`
	Position    *int   `json:"position,omitempty"`
}

type AddToCollectionRequest struct {
	BookmarkIDs []string `json:"bookmarkIds" binding:"required"`
}

type ShareBookmarkRequest struct {
	BookmarkID string `json:"bookmarkId" binding:"required"`
	SharedWith string `json:"sharedWith" binding:"required"`
	Message    string `json:"message,omitempty"`
}

type CreateNoteRequest struct {
	Content string `json:"content" binding:"required"`
	Color   string `json:"color,omitempty"`
}

type UpdateNoteRequest struct {
	Content  string `json:"content,omitempty"`
	Color    string `json:"color,omitempty"`
	IsPinned *bool  `json:"isPinned,omitempty"`
}

type CreateReminderRequest struct {
	RemindAt string `json:"remindAt" binding:"required"`
	Message  string `json:"message,omitempty"`
}

type CreatePreviewRequest struct {
	BookmarkID string `json:"bookmarkId" binding:"required"`
	URL        string `json:"url" binding:"required"`
}

type AddReadLaterRequest struct {
	BookmarkID string `json:"bookmarkId" binding:"required"`
	Priority   int    `json:"priority"`
}

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID string `json:"parentId,omitempty"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type SetExpirationRequest struct {
	ExpiresAt string `json:"expiresAt" binding:"required"`
	Action    string `json:"action,omitempty"`
}

type CreateTemplateRequest struct {
	UserID      string `json:"userId" binding:"required"`
	WorkspaceID string `json:"workspaceId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	FolderID    string `json:"folderId,omitempty"`
	TagIDs      string `json:"tagIds,omitempty"`
	Metadata    string `json:"metadata,omitempty"`
}

type ApplyTemplateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description,omitempty"`
	TargetID    string `json:"targetId" binding:"required"`
	TargetURL   string `json:"targetUrl,omitempty"`
}

// Bulk Operation DTOs

type BulkDeleteRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

type BulkMoveRequest struct {
	IDs      []string `json:"ids" binding:"required"`
	FolderID string   `json:"folderId,omitempty"`
}

type BulkTagRequest struct {
	IDs    []string `json:"ids" binding:"required"`
	TagIDs []string `json:"tagIds" binding:"required"`
}

type ReorderRequest struct {
	Items []ReorderItem `json:"items" binding:"required"`
}

type ReorderItem struct {
	ID       string `json:"id" binding:"required"`
	Position int    `json:"position"`
}
