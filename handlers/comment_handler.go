package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/test/blog/config"
	"github.com/test/blog/models"
	"github.com/test/blog/utils"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.LogError("invalid post id for comment", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid post id",
		})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("create comment validation error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError("create comment unauthorized", nil)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.GetDB().First(&post, postID).Error; err != nil {
		utils.LogError("create comment post not found", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Post not found",
		})
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}
	if err := config.GetDB().Create(&comment).Error; err != nil {
		utils.LogError("create comment database error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create comment",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Comment created successfully",
		"data": gin.H{
			"id":         comment.ID,
			"content":    comment.Content,
			"user_id":    comment.UserID,
			"username":   c.GetString("username"),
			"post_id":    comment.PostID,
			"created_at": comment.CreatedAt,
		},
	})
}

// GetComments 获取评论列表
func GetComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.LogError("get comments invalid post id", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid post id",
		})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var comments []models.Comment
	var total int64
	if err := config.GetDB().Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		utils.LogError("get comments count error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get comments count",
		})
		return
	}

	offset := (page - 1) * limit
	if err := config.GetDB().Where("post_id = ?", postID).Order("created_at desc").Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		utils.LogError("get comments list error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get comments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comments retrieved successfully",
		"data": gin.H{
			"comments": comments,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}
