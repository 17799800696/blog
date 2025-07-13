package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/test/blog/config"
	"github.com/test/blog/models"
	"github.com/test/blog/utils"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("create post validation error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError("create post unauthorized", nil)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}
	if err := config.GetDB().Create(&post).Error; err != nil {
		utils.LogError("create post database error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create post",
		})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		utils.LogError("get user for post creation", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Post created successfully",
		"data": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"username":   user.Username,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		},
	})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.LogError("update post invalid id", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid post id",
		})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("update post validation error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError("update post unauthorized", nil)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var post models.Post
	if err := config.GetDB().First(&post, postID).Error; err != nil {
		utils.LogError("update post not found", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Post not found",
		})
		return
	}

	if post.UserID != userID.(uint) {
		utils.LogError("update post forbidden", nil)
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You are not the author of this post",
		})
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	post.UpdatedAt = time.Now()
	if err := config.GetDB().Save(&post).Error; err != nil {
		utils.LogError("update post database error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post updated successfully",
		"data": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"username":   c.GetString("username"),
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		},
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.LogError("delete post invalid id", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid post id",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError("delete post unauthorized", nil)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var post models.Post
	if err := config.GetDB().First(&post, postID).Error; err != nil {
		utils.LogError("delete post not found", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Post not found",
		})
		return
	}

	if post.UserID != userID.(uint) {
		utils.LogError("delete post forbidden", nil)
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You are not the author of this post",
		})
		return
	}

	if err := config.GetDB().Delete(&post).Error; err != nil {
		utils.LogError("delete post database error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post deleted successfully",
	})
}

// GetPosts 获取文章列表
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64
	if err := config.GetDB().Model(&models.Post{}).Count(&total).Error; err != nil {
		utils.LogError("get posts count error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get posts count",
		})
		return
	}

	if err := config.GetDB().Preload("User").Order("created_at desc").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		utils.LogError("get posts list error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get posts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Posts retrieved successfully",
		"data": gin.H{
			"posts": posts,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetPost 获取单个文章
func GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.LogError("get post invalid id", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid post id",
		})
		return
	}

	var post models.Post
	if err := config.GetDB().Preload("User").First(&post, postID).Error; err != nil {
		utils.LogError("get post not found", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post retrieved successfully",
		"data":    post,
	})
}
