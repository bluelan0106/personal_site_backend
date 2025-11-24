package models

import (
    "time"

    "gorm.io/gorm"
)

// PostStatus defines post status
type PostStatus string

const (
    PostStatusDraft     PostStatus = "draft"
    PostStatusPublished PostStatus = "published"
    PostStatusArchived  PostStatus = "archived"
)

// PostVisibility defines visibility
type PostVisibility string

const (
    PostVisibilityPublic  PostVisibility = "public"
    PostVisibilityPrivate PostVisibility = "private"
)

// Post model
type Post struct {
    gorm.Model
    AuthorID    uint           `gorm:"not null;index"`
    Author      User           `gorm:"foreignKey:AuthorID"`
    Title       string         `gorm:"size:255;not null"`
    Content     string         `gorm:"type:text;not null"`
    Summary     string         `gorm:"size:500"`
    CoverImage  string         `gorm:"size:500"`
    Status      PostStatus     `gorm:"size:20;not null;default:'draft'"`
    Visibility  PostVisibility `gorm:"size:20;not null;default:'public'"`
    ViewCount   int            `gorm:"default:0"`
    PublishedAt *time.Time     `gorm:"index"`
    Tags        []Tag          `gorm:"many2many:post_tags;"`
    Comments    []Comment      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
    Reactions   []Reaction     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

// Tag model
type Tag struct {
    gorm.Model
    Name  string `gorm:"size:50;not null;uniqueIndex"`
    Slug  string `gorm:"size:50;not null;uniqueIndex"`
    Color string `gorm:"size:7;default:'#3B82F6'"`
    Posts []Post `gorm:"many2many:post_tags;"`
}

func (Post) TableName() string { return "posts" }

// BeforeSave sets PublishedAt when publishing
func (p *Post) BeforeSave(tx *gorm.DB) error {
    if p.Status == PostStatusPublished && p.PublishedAt == nil {
        now := time.Now()
        p.PublishedAt = &now
    }
    return nil
}
package models

import (
	"time"

	"gorm.io/gorm"
)

// PostStatus 定義文章狀態
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"     // 草稿
	PostStatusPublished PostStatus = "published" // 已發布
	PostStatusArchived  PostStatus = "archived"  // 封存
)

// PostVisibility 定義文章可見性
type PostVisibility string

const (
	PostVisibilityPublic  PostVisibility = "public"  // 公開
	PostVisibilityPrivate PostVisibility = "private" // 私人
)

// Post 文章模型
type Post struct {
	gorm.Model
	AuthorID   uint           `gorm:"not null;index"`                        // 作者 ID
	Author     User           `gorm:"foreignKey:AuthorID"`                   // 作者關聯
	Title      string         `gorm:"size:255;not null"`                     // 標題
	Content    string         `gorm:"type:text;not null"`                    // 內容
	Summary    string         `gorm:"size:500"`                              // 摘要
	CoverImage string         `gorm:"size:500"`                              // 封面圖
	Status     PostStatus     `gorm:"size:20;not null;default:'draft'"`      // 狀態
	Visibility PostVisibility `gorm:"size:20;not null;default:'public'"`     // 可見性
	ViewCount  int            `gorm:"default:0"`                             // 瀏覽次數
	PublishedAt *time.Time    `gorm:"index"`                                 // 發布時間
	Tags       []Tag          `gorm:"many2many:post_tags;"`                  // 標籤（多對多）
	Comments   []Comment      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"` // 留言
	Reactions  []Reaction     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"` // 反應
}

// Tag 標籤模型
type Tag struct {
	gorm.Model
	Name  string `gorm:"size:50;not null;uniqueIndex"` // 標籤名稱
	Slug  string `gorm:"size:50;not null;uniqueIndex"` // URL 友善名稱
	Color string `gorm:"size:7;default:'#3B82F6'"`     // 標籤顏色
	Posts []Post `gorm:"many2many:post_tags;"`         // 文章（多對多）
}

// BeforeSave 驗證資料
func (p *Post) BeforeSave(tx *gorm.DB) error {
	// 如果狀態變更為已發布且 PublishedAt 為空，設定發布時間
	if p.Status == PostStatusPublished && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}
