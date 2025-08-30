package task_04

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取所有文章列表
func GetPosts(c *gin.Context) {
	var posts []Post

	//预加载用户、评论信息
	queryPre := BlogDB.Preload("User").Preload("Comments")

	//分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	//查询
	if err := queryPre.Offset(offset).Limit(limit).Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "文章查询失败",
		})
		return
	}

	//查询文章总数
	var total int64
	BlogDB.Model(&Post{}).Count(&total)

	//返回信息
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "文章查询成功",
		Data: gin.H{
			"posts": posts,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})

}

// 获取单个文章详情
func GetPost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的文章ID",
		})
		return
	}

	var post Post
	if err := BlogDB.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "文章查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "文章查询成功",
		Data:    post,
	})

}

// 创建新文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的请求数据：" + err.Error(),
		})
		return
	}

	userID := GetUserId(c)
	post := Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := BlogDB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "文章创建失败",
		})
		return
	}

	//重新查询获取用户信息
	BlogDB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "文章创建成功",
		Data:    post,
	})

}

// 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的文章ID",
		})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的请求参数： " + err.Error(),
		})
		return
	}

	//查找文章
	var post Post
	if err = BlogDB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "未找到要更新的文章数据",
		})
		return
	}

	//权限校验，只有作者才能修改本人的文章
	userID := GetUserId(c)
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "只能修改自己创建的文章",
		})
		return
	}

	//更新文章
	updateParam := make(map[string]interface{})
	if req.Title != "" {
		updateParam["title"] = req.Title
	}
	if req.Content != "" {
		updateParam["content"] = req.Content
	}
	if err := BlogDB.Model(&post).Updates(updateParam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "文章更新失败",
		})
		return
	}

	//重新查询，获取完整返回信息
	BlogDB.Preload("User").First(&post, postID)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "文章更新成功",
		Data:    post,
	})

}

// 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的文章ID",
		})
		return
	}

	//查询文章信息
	var post Post
	if err := BlogDB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "未找到要更新的文章数据",
		})
		return
	}

	//校验，只能删除自己创建的文章
	userID := GetUserId(c)
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, APIResponse{
			Success: false,
			Error:   "只能修改自己创建的文章",
		})
		return
	}

	//删除文章（同时会删除相关的评论）
	if err := BlogDB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "文章删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "文章删除成功",
	})

}
