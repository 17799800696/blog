package handlers

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	PostID    uint   `json:"post_id"`
	CreatedAt string `json:"created_at"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}
