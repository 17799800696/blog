package models

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Password  string         `json:"-" gorm:"not null;size:255"` // json:"-" 表示不序列化密码字段
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100"`
	// 关联关系
	Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:UserID"`
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title     string         `json:"title" gorm:"not null;size:200"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	// 关联关系
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	PostID    uint           `json:"post_id" gorm:"not null;index"`
	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Post Post `json:"post,omitempty" gorm:"foreignKey:PostID"`
}
