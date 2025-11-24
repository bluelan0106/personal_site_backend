package post

import (
    "net/http"
    "personal_site/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type CreateCommentRequest struct {
    Content  string `json:"content" binding:"required"`
    ParentID *uint  `json:"parent_id"`
}

func CreateComment(c *gin.Context, db *gorm.DB) {
    postID := c.Param("id")
    var req CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    userID, exists := c.Get("user_id")
    if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}); return }
    var post models.Post
    if err := db.First(&post, postID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"}); return }
    comment := models.Comment{PostID: post.ID, AuthorID: userID.(uint), Content: req.Content, ParentID: req.ParentID}
    if err := db.Create(&comment).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"}); return }
    db.Preload("Author").First(&comment, comment.ID)
    c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully", "comment": comment})
}

func GetComments(c *gin.Context, db *gorm.DB) {
    postID := c.Param("id")
    var comments []models.Comment
    if err := db.Preload("Author").Preload("Replies").Preload("Reactions").Where("post_id = ? AND parent_id IS NULL AND is_deleted = ?", postID, false).Find(&comments).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"}); return }
    c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func UpdateComment(c *gin.Context, db *gorm.DB) {
    commentID := c.Param("comment_id")
    var req struct{ Content string `json:"content" binding:"required"` }
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    userID, exists := c.Get("user_id")
    if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}); return }
    var comment models.Comment
    if err := db.First(&comment, commentID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"}); return }
    if comment.AuthorID != userID.(uint) { c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own comments"}); return }
    comment.Content = req.Content
    comment.IsEdited = true
    if err := db.Save(&comment).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"}); return }
    c.JSON(http.StatusOK, gin.H{"message": "Comment updated", "comment": comment})
}

func DeleteComment(c *gin.Context, db *gorm.DB) {
    commentID := c.Param("comment_id")
    userID, exists := c.Get("user_id")
    if !exists { c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}); return }
    var comment models.Comment
    if err := db.First(&comment, commentID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"}); return }
    if comment.AuthorID != userID.(uint) { c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"}); return }
    comment.IsDeleted = true
    comment.Content = "[此留言已刪除]"
    if err := db.Save(&comment).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"}); return }
    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
package post

import (
	"net/http"
	"personal_site/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCommentRequest 建立留言請求
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	ParentID *uint  `json:"parent_id"` // 如果是回覆，提供父留言 ID
}

// UpdateCommentRequest 更新留言請求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// CreateComment 建立留言
func CreateComment(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")
	var req CreateCommentRequest
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

	// 如果是回覆，檢查父留言是否存在
	if req.ParentID != nil {
		var parentComment models.Comment
		if err := db.First(&parentComment, *req.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
			return
		}
		// 確保父留言屬於同一篇文章
		if parentComment.PostID != post.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment does not belong to this post"})
			return
		}
	}

	// 建立留言
	comment := models.Comment{
		PostID:   post.ID,
		AuthorID: userID.(uint),
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 載入作者資料
	db.Preload("Author").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetComments 取得文章的所有留言
func GetComments(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")

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

	var comments []models.Comment
	// 只取得頂層留言（沒有父留言的）
	if err := db.Where("post_id = ? AND parent_id IS NULL AND is_deleted = ?", postID, false).
		Preload("Author").
		Preload("Reactions").
		Preload("Replies", "is_deleted = ?", false).
		Preload("Replies.Author").
		Preload("Replies.Reactions").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// UpdateComment 更新留言
func UpdateComment(c *gin.Context, db *gorm.DB) {
	commentID := c.Param("comment_id")
	var req UpdateCommentRequest
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

	// 檢查是否為作者
	if comment.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own comments"})
		return
	}

	// 更新留言
	if err := db.Model(&comment).Updates(map[string]interface{}{
		"content":   req.Content,
		"is_edited": true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// 重新載入資料
	db.Preload("Author").First(&comment, comment.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

// DeleteComment 刪除留言（軟刪除）
func DeleteComment(c *gin.Context, db *gorm.DB) {
	commentID := c.Param("comment_id")

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

	// 檢查權限
	role, _ := c.Get("role")
	if comment.AuthorID != userID.(uint) && role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
		return
	}

	// 軟刪除（標記為已刪除，保留結構）
	if err := db.Model(&comment).Updates(map[string]interface{}{
		"is_deleted": true,
		"content":    "[此留言已刪除]",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
