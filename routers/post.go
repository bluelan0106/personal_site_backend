package routers

import (
    "personal_site/controllers/post"
    "personal_site/middlewares"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type postRouter struct{}

func (postRouter) RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
    posts := r.Group("/posts")
    {
        posts.GET("", middlewares.AuthOptional(), func(c *gin.Context) { post.GetPosts(c, db) })
        posts.GET(":id", middlewares.AuthOptional(), func(c *gin.Context) { post.GetPost(c, db) })
        posts.GET(":id/reactions", middlewares.AuthOptional(), func(c *gin.Context) { post.GetPostReactionsSummary(c, db) })

        posts.POST("", middlewares.AuthRequired(), func(c *gin.Context) { post.CreatePost(c, db) })
        posts.PUT(":id", middlewares.AuthRequired(), func(c *gin.Context) { post.UpdatePost(c, db) })
        posts.DELETE(":id", middlewares.AuthRequired(), func(c *gin.Context) { post.DeletePost(c, db) })

        posts.POST(":id/reactions", middlewares.AuthRequired(), func(c *gin.Context) { post.AddReactionToPost(c, db) })
        posts.GET(":id/comments", middlewares.AuthOptional(), func(c *gin.Context) { post.GetComments(c, db) })
        posts.POST(":id/comments", middlewares.AuthRequired(), func(c *gin.Context) { post.CreateComment(c, db) })
    }

    comments := r.Group("/comments")
    comments.Use(middlewares.AuthRequired())
    {
        comments.PUT(":comment_id", func(c *gin.Context) { post.UpdateComment(c, db) })
        comments.DELETE(":comment_id", func(c *gin.Context) { post.DeleteComment(c, db) })
        comments.POST(":comment_id/reactions", func(c *gin.Context) { post.AddReactionToComment(c, db) })
        comments.GET(":comment_id/reactions", middlewares.AuthOptional(), func(c *gin.Context) { post.GetCommentReactionsSummary(c, db) })
    }
}
 
