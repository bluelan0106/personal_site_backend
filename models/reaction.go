package models

import (
    "gorm.io/gorm"
)

// ReactionType
type ReactionType string

const (
    ReactionTypeLike  ReactionType = "like"
    ReactionTypeLove  ReactionType = "love"
    ReactionTypeHaha  ReactionType = "haha"
    ReactionTypeWow   ReactionType = "wow"
    ReactionTypeSad   ReactionType = "sad"
    ReactionTypeAngry ReactionType = "angry"
    ReactionTypeCare  ReactionType = "care"
)

// Reaction model (polymorphic between Post and Comment)
type Reaction struct {
    gorm.Model
    UserID    uint         `gorm:"not null;index:idx_reaction_unique,unique"`
    User      User         `gorm:"foreignKey:UserID"`
    PostID    *uint        `gorm:"index:idx_reaction_unique,unique"`
    Post      *Post        `gorm:"foreignKey:PostID"`
    CommentID *uint        `gorm:"index:idx_reaction_unique,unique"`
    Comment   *Comment     `gorm:"foreignKey:CommentID"`
    Type      ReactionType `gorm:"size:20;not null"`
}

func (Reaction) TableName() string { return "reactions" }

func (r *Reaction) BeforeSave(tx *gorm.DB) error {
    if (r.PostID == nil && r.CommentID == nil) || (r.PostID != nil && r.CommentID != nil) {
        return gorm.ErrInvalidData
    }
    return nil
}
package models

import (
	"gorm.io/gorm"
)

// ReactionType 反應類型
type ReactionType string

const (
	ReactionTypeLike    ReactionType = "like"    // 讚
	ReactionTypeLove    ReactionType = "love"    // 愛心
	ReactionTypeHaha    ReactionType = "haha"    // 哈哈
	ReactionTypeWow     ReactionType = "wow"     // 驚訝
	ReactionTypeSad     ReactionType = "sad"     // 難過
	ReactionTypeAngry   ReactionType = "angry"   // 生氣
	ReactionTypeCare    ReactionType = "care"    // 關心
)

// Reaction 反應模型（用於文章和留言）
type Reaction struct {
	gorm.Model
	UserID    uint         `gorm:"not null;index:idx_reaction_unique,unique"` // 使用者 ID
	User      User         `gorm:"foreignKey:UserID"`                         // 使用者關聯
	PostID    *uint        `gorm:"index:idx_reaction_unique,unique"`          // 文章 ID（可為空）
	Post      *Post        `gorm:"foreignKey:PostID"`                         // 文章關聯
	CommentID *uint        `gorm:"index:idx_reaction_unique,unique"`          // 留言 ID（可為空）
	Comment   *Comment     `gorm:"foreignKey:CommentID"`                      // 留言關聯
	Type      ReactionType `gorm:"size:20;not null"`                          // 反應類型
}

// TableName 指定表名
func (Reaction) TableName() string {
	return "reactions"
}

// BeforeSave 驗證反應必須綁定到文章或留言
func (r *Reaction) BeforeSave(tx *gorm.DB) error {
	// 確保反應只能對文章或留言其中一個
	if (r.PostID == nil && r.CommentID == nil) || (r.PostID != nil && r.CommentID != nil) {
		return gorm.ErrInvalidData
	}
	return nil
}
