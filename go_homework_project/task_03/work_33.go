package task_03

import (
	"fmt"
	"go_homework_project/task_03/repository"
	"strconv"

	"gorm.io/gorm"
)

/*
查询某个用户发布的所有文章及其对应的评论信息
*/
func GetAllInfoByUserName(name string) User {
	var user User

	if err := repository.DB.Preload("Posts.Comments").Where("name = ?", name).First(&user).Error; err != nil {
		fmt.Println("查询失败：", err)
	} else {
		fmt.Printf("作者：%v\n", user.Name)
		for _, post := range user.Posts {
			fmt.Printf("文章题目：%s，文章内容：%s\n", post.Title, post.Content)
			for _, comment := range post.Comments {
				fmt.Printf("评论：%s\n", comment.Content)
			}
		}
	}

	return user
}

/*
查询评论数量最多的文章信息。

先查出最大评论数，然后根据最大评论数查询相关文章信息
*/

// 最大评论数文章信息结构体
type HottestPost struct {
	ID         uint
	Title      string
	Author     string
	CommentNum uint
}

func GetHottestPostInfo() []HottestPost {
	var maxCommentNum int

	//查询最大评论数量
	repository.DB.Table("comments").
		Select("count(*) as comment_num").
		Group("post_id").
		Order("comment_num desc").
		Limit(1).
		Scan(&maxCommentNum)

	if maxCommentNum == 0 {
		fmt.Println("所有文章均无评论！")
		return nil
	}

	//根据最大评论数查询相关文章信息
	var hottestPosts []HottestPost
	result := repository.DB.Table("posts").
		Select("posts.id,posts.title, users.name as author, "+strconv.Itoa(maxCommentNum)+" as comment_num").
		Joins("join users on users.id = posts.user_id").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id,users.name").
		Having("COUNT(comments.id) = ?", maxCommentNum).
		Find(&hottestPosts)
	if result.Error != nil {
		fmt.Println("查询失败：", result.Error)
		return nil
	}

	return hottestPosts
}

/*
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
*/
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	if err := tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_num", gorm.Expr("(SELECT COUNT(id) FROM posts WHERE user_id = ?)", p.UserID)).Error; err != nil {
		return err
	}
	return nil
}

func AddPost(posts []Post) {
	if err := repository.DB.Create(&posts).Error; err != nil {
		fmt.Println("新增数据失败：", err)
	}
}

func AddUserAndPost(users []User) {
	if err := repository.DB.Create(&users).Error; err != nil {
		fmt.Println("新增数据失败：", err)
	}
}

/*
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	var post_id = c.PostID
	var commentNum int
	var comment_status int
	if err := tx.Model(&Comment{}).Where("post_id = ?", post_id).Select("count(*) as commentNum").Scan(&commentNum).Error; err != nil {
		fmt.Println("查询commentNum失败：", err)
		return err
	}

	if commentNum <= 1 {
		comment_status = 0
	} else {
		comment_status = 1
	}

	if err := tx.Model(&Post{}).Where("id = ?", post_id).Update("comment_status", comment_status).Error; err != nil {
		fmt.Println("更新comment_status失败：", err)
		return err
	}

	return nil
}

func DeleteCommentsById(ids []uint) {

	for _, id := range ids {
		comment := Comment{}
		if err := repository.DB.First(&comment, id).Error; err != nil {
			fmt.Println("删除评论失败：", err)
		}
		if err := repository.DB.Delete(&comment).Error; err != nil {
			fmt.Println("删除评论失败：", err)
		}
	}
}
