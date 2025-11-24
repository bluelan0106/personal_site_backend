package post

import (
	"net/http"
	"personal_site/models"
	"strconv"
	"time"
    
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatePostRequest struct {
	Title      string                `json:"title" binding:"required,min=1,max=255"`
	Content    string                `json:"content" binding:"required"`
	Summary    string                `json:"summary" binding:"max=500"`
	CoverImage string                `json:"cover_image" binding:"omitempty,url"`
	Status     models.PostStatus     `json:"status" binding:"omitempty,oneof=draft published archived"`
	Visibility models.PostVisibility `json:"visibility" binding:"omitempty,oneof=public private"`
	Tags       []string              `json:"tags"`
}

type UpdatePostRequest struct {
	Title      *string                `json:"title" binding:"omitempty,min=1,max=255"`
	Content    *string                `json:"content" binding:"omitempty"`
	Summary    *string                `json:"summary" binding:"omitempty,max=500"`
	CoverImage *string                `json:"cover_image" binding:"omitempty,url"`
	Status     *models.PostStatus     `json:"status" binding:"omitempty,oneof=draft published archived"`
	Visibility *models.PostVisibility `json:"visibility" binding:"omitempty,oneof=public private"`
	Tags       []string               `json:"tags"`
}

func CreatePost(c *gin.Context, db *gorm.DB) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if req.Status == "" {
		req.Status = models.PostStatusDraft
	}
	if req.Visibility == "" {
		req.Visibility = models.PostVisibilityPublic
	}
	post := models.Post{
		AuthorID:   userID.(uint),
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
		Status:     req.Status,
		Visibility: req.Visibility,
	}
	tx := db.Begin()
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	if len(req.Tags) > 0 {
		var tags []models.Tag
		for _, tagName := range req.Tags {
			var tag models.Tag
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName, Slug: tagName}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tags"})
				return
			}
			tags = append(tags, tag)
		}
		if err := tx.Model(&post).Association("Tags").Append(tags); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags"})
			return
		}
	}
	tx.Commit()
	db.Preload("Author").Preload("Tags").First(&post, post.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": post})
}

func GetPosts(c *gin.Context, db *gorm.DB) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	tag := c.Query("tag")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	query := db.Model(&models.Post{}).Preload("Author").Preload("Tags").Preload("Reactions")
	userID, exists := c.Get("user_id")
	if !exists {
		query = query.Where("status = ? AND visibility = ?", models.PostStatusPublished, models.PostVisibilityPublic)
	} else {
		query = query.Where("(author_id = ? OR (status = ? AND visibility = ?))", userID.(uint), models.PostStatusPublished, models.PostVisibilityPublic)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if tag != "" {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").Joins("JOIN tags ON tags.id = post_tags.tag_id").Where("tags.name = ?", tag)
	}
	var total int64
	query.Count(&total)
	var posts []models.Post
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "pagination": gin.H{"page": page, "page_size": pageSize, "total": total, "total_pages": (total + int64(pageSize) - 1) / int64(pageSize)}})
}

func GetPost(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")
	var post models.Post
	query := db.Preload("Author").Preload("Tags").Preload("Comments", "is_deleted = ?", false).Preload("Comments.Author").Preload("Comments.Reactions").Preload("Reactions").Preload("Reactions.User")
	if err := query.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}
	userID, exists := c.Get("user_id")
	if post.Visibility == models.PostVisibilityPrivate {
		if !exists || userID.(uint) != post.AuthorID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}
	db.Model(&post).Update("view_count", gorm.Expr("view_count + 1"))
	post.ViewCount++
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func UpdatePost(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")
	var req UpdatePostRequest
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
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own posts"})
		return
	}
	tx := db.Begin()
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.CoverImage != nil {
		updates["cover_image"] = *req.CoverImage
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		if *req.Status == models.PostStatusPublished && post.PublishedAt == nil {
			now := time.Now()
			updates["published_at"] = now
		}
	}
	if req.Visibility != nil {
		updates["visibility"] = *req.Visibility
	}
	if len(updates) > 0 {
		if err := tx.Model(&post).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
			return
		}
	}
	if req.Tags != nil {
		var tags []models.Tag
		for _, tagName := range req.Tags {
			var tag models.Tag
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName, Slug: tagName}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tags"})
				return
			}
			tags = append(tags, tag)
		}
		if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags"})
			return
		}
	}
	tx.Commit()
	db.Preload("Author").Preload("Tags").First(&post, post.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Post updated", "post": post})
}

func DeletePost(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")
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
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}
	if err := db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
 
