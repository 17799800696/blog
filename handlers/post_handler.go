package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/test/blog/config"
	"github.com/test/blog/models"
	"github.com/test/blog/utils"
	"go.uber.org/zap"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(err, "create post validation error")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError(utils.ErrUnauthorized, "create post unauthorized")
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// 创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := config.DB.Create(&post).Error; err != nil {
		utils.LogError(err, "create post database error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to create post",
		})
		return
	}

	// 获取用户信息
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.LogError(err, "get user for post creation")
	}

	utils.LogInfo("post created successfully",
		zap.Uint("post_id", post.ID),
		zap.Uint("user_id", userID.(uint)),
	)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Post created successfully",
		Data: PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			UserID:    post.UserID,
			Username:  user.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// GetPosts 获取文章列表
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64

	// 获取总数
	if err := config.DB.Model(&models.Post{}).Count(&total).Error; err != nil {
		utils.LogError(err, "get posts count error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to get posts count",
		})
		return
	}

	// 获取文章列表（包含用户信息）
	if err := config.DB.Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error; err != nil {
		utils.LogError(err, "get posts database error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to get posts",
		})
		return
	}

	// 构建响应
	var postResponses []PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			UserID:    post.UserID,
			Username:  post.User.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		})
	}

	utils.LogInfo("posts retrieved successfully",
		zap.Int("count", len(postResponses)),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Posts retrieved successfully",
		Data: PostListResponse{
			Posts: postResponses,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

// GetPost 获取单个文章
func GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError(err, "invalid post id")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var post models.Post
	if err := config.DB.Preload("User").First(&post, postID).Error; err != nil {
		utils.LogError(err, "get post not found")
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	utils.LogInfo("post retrieved successfully",
		zap.Uint("post_id", post.ID),
	)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Post retrieved successfully",
		Data: PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			UserID:    post.UserID,
			Username:  post.User.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError(err, "invalid post id for update")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(err, "update post validation error")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError(utils.ErrUnauthorized, "update post unauthorized")
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.LogError(err, "update post not found")
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.LogError(utils.ErrForbidden, "update post forbidden")
		c.JSON(http.StatusForbidden, Response{
			Success: false,
			Message: "You can only update your own posts",
		})
		return
	}

	// 更新文章
	updates := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}

	if err := config.DB.Model(&post).Updates(updates).Error; err != nil {
		utils.LogError(err, "update post database error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to update post",
		})
		return
	}

	// 重新获取文章信息
	config.DB.Preload("User").First(&post, postID)

	utils.LogInfo("post updated successfully",
		zap.Uint("post_id", post.ID),
		zap.Uint("user_id", userID.(uint)),
	)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Post updated successfully",
		Data: PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			UserID:    post.UserID,
			Username:  post.User.Username,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError(err, "invalid post id for delete")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError(utils.ErrUnauthorized, "delete post unauthorized")
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.LogError(err, "delete post not found")
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		utils.LogError(utils.ErrForbidden, "delete post forbidden")
		c.JSON(http.StatusForbidden, Response{
			Success: false,
			Message: "You can only delete your own posts",
		})
		return
	}

	// 删除文章（软删除）
	if err := config.DB.Delete(&post).Error; err != nil {
		utils.LogError(err, "delete post database error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to delete post",
		})
		return
	}

	utils.LogInfo("post deleted successfully",
		zap.Uint("post_id", post.ID),
		zap.Uint("user_id", userID.(uint)),
	)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Post deleted successfully",
	})
}
