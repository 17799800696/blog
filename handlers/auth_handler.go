package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/test/blog/config"
	"github.com/test/blog/models"
	"github.com/test/blog/utils"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, Response{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, Response{
			Success: false,
			Message: "Email already exists",
		})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to hash password",
		})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	// 返回响应
	authResponse := AuthResponse{
		Token: token,
		User: struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "User registered successfully",
		Data:    authResponse,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	// 返回响应
	authResponse := AuthResponse{
		Token: token,
		User: struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Login successful",
		Data:    authResponse,
	})
}

// GetProfile 获取用户信息
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Profile retrieved successfully",
		Data: gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
