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

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError(err, "invalid post id for comment")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(err, "create comment validation error")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError(utils.ErrUnauthorized, "create comment unauthorized")
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.LogError(err, "create comment post not found")
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		utils.LogError(err, "create comment database error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to create comment",
		})
		return
	}

	// 获取用户信息
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.LogError(err, "get user for comment creation")
	}

	utils.LogInfo("comment created successfully",
		zap.Uint("comment_id", comment.ID),
		zap.Uint("post_id", uint(postID)),
		zap.Uint("user_id", userID.(uint)),
	)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Comment created successfully",
		Data: CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			Username:  user.Username,
			PostID:    comment.PostID,
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		},
	})
}

// GetComments 获取文章评论列表
func GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError(err, "invalid post id for get comments")
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		utils.LogError(err, "get comments post not found")
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var comments []models.Comment
	var total int64

	// 获取评论总数
	if err := config.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		utils.LogError(err, "get comments count error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to get comments count",
		})
		return
	}

	// 获取评论列表（包含用户信息）
	if err := config.DB.Preload("User").Where("post_id = ?", postID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error; err != nil {
		utils.LogError(err, "get comments database error")
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to get comments",
		})
		return
	}

	// 构建响应
	var commentResponses []CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			Username:  comment.User.Username,
			PostID:    comment.PostID,
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		})
	}

	utils.LogInfo("comments retrieved successfully",
		zap.Uint("post_id", uint(postID)),
		zap.Int("count", len(commentResponses)),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Comments retrieved successfully",
		Data: CommentListResponse{
			Comments: commentResponses,
			Total:    total,
			Page:     page,
			Limit:    limit,
		},
	})
}
