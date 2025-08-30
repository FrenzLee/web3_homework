package task_04

import "gorm.io/gorm"

// 用户信息模型
type User struct {
	gorm.Model
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"-"` //请求查询时，此字段不给前端返回
	Email    string    `gorm:"unique;not null" json:"email"`
	Posts    []Post    `json:"posts,omitempty"` //omitempty-没有相关信息，就不给前端返回
	Comments []Comment `json:"comments,omitempty"`
}

// 文章模型
type Post struct {
	gorm.Model
	Title    string    `gorm:"not null" json:"title"`
	Content  string    `gorm:"not null" json:"content"`
	UserID   uint      `json:"user_id"`
	User     User      `json:"user,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user,omitempty"`
	PostID  uint   `json:"post_id"`
	Post    Post   `json:"post,omitempty"`
}

// 账号注册入参
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// 用户登录入参
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

// 接口返回模型
type APIResponse struct {
	Success bool        `json:"succcess"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
