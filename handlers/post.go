package handlers

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

// PostResponse 文章响应
type PostResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PostListResponse 文章列表响应
type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
