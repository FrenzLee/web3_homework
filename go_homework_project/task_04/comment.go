package task_04

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取指定文章的所有评论
func GetCommentsByPost(c *gin.Context) {
	postIdParam := c.Param("postId")
	postID, err := strconv.ParseUint(postIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的文章ID",
		})
		return
	}

	var post Post
	if err = BlogDB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "文章数据不存在",
		})
		return
	}

	//预加载用户信息
	queryPre := BlogDB.Where("post_id = ?", postID).Preload("User")

	//分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	//查询评论信息
	var comments []Comment
	if err := queryPre.Offset(offset).Limit(limit).Order("created_at desc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "评论查询失败",
		})
		return
	}

	//查询评论总数
	var total int64
	BlogDB.Model(&Comment{}).Where("post_id = ?", postID).Count(&total)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "评论查询成功",
		Data: gin.H{
			"comments": comments,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})

}

// 创建新评论
func CreateCommentForPost(c *gin.Context) {
	var reqCom CreateCommentRequest
	if err := c.ShouldBindJSON(&reqCom); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的请求数据：" + err.Error(),
		})
		return
	}

	//查找文章
	var post Post
	if err := BlogDB.First(&post, reqCom.PostID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "未找到要评论的文章数据",
		})
		return
	}

	userID := GetUserId(c)
	comment := Comment{
		Content: reqCom.Content,
		UserID:  userID,
		PostID:  reqCom.PostID,
	}

	//创建评论
	if err := BlogDB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "评论创建失败",
		})
		return
	}

	//重新查询获取用户信息
	BlogDB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "评论创建成功",
		Data:    comment,
	})

}
