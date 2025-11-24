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
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    userID, exists := c.Get("user_id")
    if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}); return }
    var post models.Post
    if err := db.First(&post, postID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"}); return }
    var reaction models.Reaction
    if err := db.Where("user_id = ? AND post_id = ?", userID.(uint), post.ID).First(&reaction).Error; err == nil {
        if reaction.Type == req.Type {
            db.Delete(&reaction)
            c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
            return
        }
        reaction.Type = req.Type
        db.Save(&reaction)
        c.JSON(http.StatusOK, gin.H{"message": "Reaction updated", "reaction": reaction})
        return
    }
    reaction = models.Reaction{UserID: userID.(uint), PostID: &post.ID, Type: req.Type}
    if err := db.Create(&reaction).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"}); return }
    c.JSON(http.StatusCreated, gin.H{"message": "Reaction added", "reaction": reaction})
}

func AddReactionToComment(c *gin.Context, db *gorm.DB) {
    commentID := c.Param("comment_id")
    var req ReactionRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    userID, exists := c.Get("user_id")
    if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}); return }
    var comment models.Comment
    if err := db.First(&comment, commentID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"}); return }
    var reaction models.Reaction
    if err := db.Where("user_id = ? AND comment_id = ?", userID.(uint), comment.ID).First(&reaction).Error; err == nil {
        if reaction.Type == req.Type { db.Delete(&reaction); c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"}); return }
        reaction.Type = req.Type
        db.Save(&reaction)
        c.JSON(http.StatusOK, gin.H{"message": "Reaction updated", "reaction": reaction})
        return
    }
    reaction = models.Reaction{UserID: userID.(uint), CommentID: &comment.ID, Type: req.Type}
    if err := db.Create(&reaction).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"}); return }
    c.JSON(http.StatusCreated, gin.H{"message": "Reaction added", "reaction": reaction})
}

func GetPostReactionsSummary(c *gin.Context, db *gorm.DB) {
    postID := c.Param("id")
    type ReactionCount struct { Type models.ReactionType `json:"type"`; Count int64 `json:"count"` }
    var reactions []ReactionCount
    if err := db.Model(&models.Reaction{}).Select("type, COUNT(*) as count").Where("post_id = ?", postID).Group("type").Find(&reactions).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"}); return }
    var userReaction *models.Reaction
    if userID, exists := c.Get("user_id"); exists {
        var r models.Reaction
        if err := db.Where("user_id = ? AND post_id = ?", userID.(uint), postID).First(&r).Error; err == nil { userReaction = &r }
    }
    c.JSON(http.StatusOK, gin.H{"reactions": reactions, "user_reaction": userReaction})
}

func GetCommentReactionsSummary(c *gin.Context, db *gorm.DB) {
    commentID := c.Param("comment_id")
    type ReactionCount struct { Type models.ReactionType `json:"type"`; Count int64 `json:"count"` }
    var reactions []ReactionCount
    if err := db.Model(&models.Reaction{}).Select("type, COUNT(*) as count").Where("comment_id = ?", commentID).Group("type").Find(&reactions).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"}); return }
    var userReaction *models.Reaction
    if userID, exists := c.Get("user_id"); exists {
        var r models.Reaction
        if err := db.Where("user_id = ? AND comment_id = ?", userID.(uint), commentID).First(&r).Error; err == nil { userReaction = &r }
    }
    c.JSON(http.StatusOK, gin.H{"reactions": reactions, "user_reaction": userReaction})
}
package post

import (
	"net/http"
	"personal_site/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddReactionRequest 新增反應請求
type AddReactionRequest struct {
	Type models.ReactionType `json:"type" binding:"required,oneof=like love haha wow sad angry care"`
}

// AddReactionToPost 對文章新增反應
func AddReactionToPost(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")
	var req AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 檢查文章是否存在
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	// 檢查是否已經有反應
	var existingReaction models.Reaction
	err := db.Where("user_id = ? AND post_id = ?", userID.(uint), post.ID).First(&existingReaction).Error

	if err == nil {
		// 已有反應，更新類型
		if existingReaction.Type == req.Type {
			// 相同類型，刪除反應（取消）
			if err := db.Delete(&existingReaction).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
			return
		}
		// 不同類型，更新
		if err := db.Model(&existingReaction).Update("type", req.Type).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "Reaction updated",
			"reaction": existingReaction,
		})
		return
	}

	// 沒有反應，建立新的
	reaction := models.Reaction{
		UserID: userID.(uint),
		PostID: &post.ID,
		Type:   req.Type,
	}

	if err := db.Create(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Reaction added",
		"reaction": reaction,
	})
}

// AddReactionToComment 對留言新增反應
func AddReactionToComment(c *gin.Context, db *gorm.DB) {
	commentID := c.Param("comment_id")
	var req AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 檢查留言是否存在
	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}

	// 檢查是否已經有反應
	var existingReaction models.Reaction
	err := db.Where("user_id = ? AND comment_id = ?", userID.(uint), comment.ID).First(&existingReaction).Error

	if err == nil {
		// 已有反應
		if existingReaction.Type == req.Type {
			// 相同類型，刪除（取消）
			if err := db.Delete(&existingReaction).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
			return
		}
		// 不同類型，更新
		if err := db.Model(&existingReaction).Update("type", req.Type).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "Reaction updated",
			"reaction": existingReaction,
		})
		return
	}

	// 沒有反應，建立新的
	reaction := models.Reaction{
		UserID:    userID.(uint),
		CommentID: &comment.ID,
		Type:      req.Type,
	}

	if err := db.Create(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Reaction added",
		"reaction": reaction,
	})
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
