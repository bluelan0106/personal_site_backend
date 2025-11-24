package models

import (
	"gorm.io/gorm"
)

// ReactionType defines reaction kinds
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

// Reaction model
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
