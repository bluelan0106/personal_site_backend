package models

import (
    "gorm.io/gorm"
)

// Comment model
type Comment struct {
    gorm.Model
    PostID      uint       `gorm:"not null;index"`
    Post        Post       `gorm:"foreignKey:PostID"`
    AuthorID    uint       `gorm:"not null;index"`
    Author      User       `gorm:"foreignKey:AuthorID"`
    Content     string     `gorm:"type:text;not null"`
    ParentID    *uint      `gorm:"index"`
    Parent      *Comment   `gorm:"foreignKey:ParentID"`
    Replies     []Comment  `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
    Reactions   []Reaction `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
    IsEdited    bool       `gorm:"default:false"`
    IsDeleted   bool       `gorm:"default:false;index"`
}

func (Comment) TableName() string { return "comments" }
package models

import (
	"gorm.io/gorm"
)

// Comment 留言模型
type Comment struct {
	gorm.Model
	PostID      uint       `gorm:"not null;index"`                             // 文章 ID
	Post        Post       `gorm:"foreignKey:PostID"`                          // 文章關聯
	AuthorID    uint       `gorm:"not null;index"`                             // 作者 ID
	Author      User       `gorm:"foreignKey:AuthorID"`                        // 作者關聯
	Content     string     `gorm:"type:text;not null"`                         // 留言內容
	ParentID    *uint      `gorm:"index"`                                      // 父留言 ID（用於回覆）
	Parent      *Comment   `gorm:"foreignKey:ParentID"`                        // 父留言關聯
	Replies     []Comment  `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"` // 子留言
	Reactions   []Reaction `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"` // 反應
	IsEdited    bool       `gorm:"default:false"`                              // 是否已編輯
	IsDeleted   bool       `gorm:"default:false;index"`                        // 軟刪除標記
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
