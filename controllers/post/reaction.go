package post

import (
	"net/http"
	"personal_site/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReactionRequest struct {
	Type models.ReactionType `json:"type" binding:"required,oneof=like love haha wow sad angry care"`
}

func AddReactionToPost(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")
	var req ReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}
	var existing models.Reaction
	if err := db.Where("user_id = ? AND post_id = ?", userID.(uint), post.ID).First(&existing).Error; err == nil {
		if existing.Type == req.Type {
			db.Delete(&existing)
			c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
			return
		}
		existing.Type = req.Type
		db.Save(&existing)
		c.JSON(http.StatusOK, gin.H{"message": "Reaction updated", "reaction": existing})
		return
	}
	r := models.Reaction{UserID: userID.(uint), PostID: &post.ID, Type: req.Type}
	if err := db.Create(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Reaction added", "reaction": r})
}

func AddReactionToComment(c *gin.Context, db *gorm.DB) {
	commentID := c.Param("comment_id")
	var req ReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}
	var existing models.Reaction
	if err := db.Where("user_id = ? AND comment_id = ?", userID.(uint), comment.ID).First(&existing).Error; err == nil {
		if existing.Type == req.Type {
			db.Delete(&existing)
			c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
			return
		}
		existing.Type = req.Type
		db.Save(&existing)
		c.JSON(http.StatusOK, gin.H{"message": "Reaction updated", "reaction": existing})
		return
	}
	r := models.Reaction{UserID: userID.(uint), CommentID: &comment.ID, Type: req.Type}
	if err := db.Create(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Reaction added", "reaction": r})
}

// GetPostReactionsSummary 取得文章的反應統計
func GetPostReactionsSummary(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")

	type ReactionCount struct {
		Type  models.ReactionType `json:"type"`
		Count int64               `json:"count"`
	}

	var reactions []ReactionCount
	if err := db.Model(&models.Reaction{}).
		Select("type, COUNT(*) as count").
		Where("post_id = ?", postID).
		Group("type").
		Find(&reactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"})
		return
	}

	// 取得當前使用者的反應
	var userReaction *models.Reaction
	if userID, exists := c.Get("user_id"); exists {
		var reaction models.Reaction
		if err := db.Where("user_id = ? AND post_id = ?", userID.(uint), postID).First(&reaction).Error; err == nil {
			userReaction = &reaction
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"reactions":     reactions,
		"user_reaction": userReaction,
	})
}

// GetCommentReactionsSummary 取得留言的反應統計
func GetCommentReactionsSummary(c *gin.Context, db *gorm.DB) {
	commentID := c.Param("comment_id")

	type ReactionCount struct {
		Type  models.ReactionType `json:"type"`
		Count int64               `json:"count"`
	}

	var reactions []ReactionCount
	if err := db.Model(&models.Reaction{}).
		Select("type, COUNT(*) as count").
		Where("comment_id = ?", commentID).
		Group("type").
		Find(&reactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"})
		return
	}

	// 取得當前使用者的反應
	var userReaction *models.Reaction
	if userID, exists := c.Get("user_id"); exists {
		var reaction models.Reaction
		if err := db.Where("user_id = ? AND comment_id = ?", userID.(uint), commentID).First(&reaction).Error; err == nil {
			userReaction = &reaction
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"reactions":     reactions,
		"user_reaction": userReaction,
	})
}
 
