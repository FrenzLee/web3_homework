package task_03

import (
	"fmt"
	"go_homework_project/task_03/repository"

	"gorm.io/gorm"
)

/*
名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、
age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
*/

type Student struct {
	ID    uint
	Name  string
	Age   uint
	Grade string
}

/*
accounts 表（包含字段 id 主键， balance 账户余额）
*/
type Account struct {
	ID        uint `gorm:"primaryKey"`
	AccountNo string
	Balance   int64
}

/*
transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）
*/
type Transaction struct {
	ID              uint `gorm:"primaryKey"`
	From_account_id uint
	To_account_id   uint
	Amount          int64
}

/*
employees 表，包含字段 id 、 name 、 department 、 salary
*/
type Employee struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Department string
	Salary     int64
}

/*
books 表，包含字段 id 、 title 、 author 、 price
*/
type Book struct {
	ID     uint `gorm:"primaryKey"`
	Title  string
	Author string
	Price  int64
}

/*
使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
*/
type User struct {
	gorm.Model
	//ID       uint
	Name     string
	post_num uint //文章数量
	Posts    []Post
}

type Post struct {
	gorm.Model
	//ID             uint
	Title          string
	Content        string
	UserID         uint
	comment_status uint //评论状态：0-无评论，1-有评论
	Comments       []Comment
}

type Comment struct {
	//gorm.Model
	ID      uint
	Content string
	PostID  uint
}

func InitTable() {
	//repository.DB.AutoMigrate(&Student{})
	//repository.DB.AutoMigrate(&Account{}, &Transaction{})
	//repository.DB.AutoMigrate(&Employee{}, &Book{})
	repository.DB.AutoMigrate(&Post{}, &Comment{})
}

func InitData() {
	repository.DB.Create([]Account{
		{AccountNo: "A", Balance: 200},
		{AccountNo: "B", Balance: 300},
	})
}

func InitDataBySqlx() {
	employees := []Employee{
		{Name: "张三", Department: "生产部", Salary: 12000},
		{Name: "李四", Department: "技术部", Salary: 15000},
		{Name: "王五", Department: "行政部", Salary: 8000},
	}

	sql := `insert into employees (name, department, salary) 
			values (:name, :department, :salary)`

	for _, employee := range employees {
		_, err := repository.DB_Sqlx.NamedExec(sql, employee)
		if err != nil {
			fmt.Println("插入数据出错：", err)
		}
	}

	fmt.Println("插入数据完成")
}

func InitDataBySqlx1() {
	books := []Book{
		{Title: "go语言经典学习手册", Author: "Jk.Wang", Price: 49},
		{Title: "高效的go语言", Author: "Sophy.Liu", Price: 58},
		{Title: "go语言源码解析", Author: "Frank.Lee", Price: 93},
	}

	sql := `insert into books (title, author,price) values (:title, :author,:price)`

	for _, book := range books {
		_, err := repository.DB_Sqlx.NamedExec(sql, book)
		if err != nil {
			fmt.Println("插入数据出错：", err)
		}
	}

	fmt.Println("插入数据完成")
}

func InitDataBlog() {
	/*newUsers := []User{
		{Name: "作者1",
			Posts: []Post{
				{Title: "作者1的第一篇文章", Content: "静夜思", Comments: []Comment{{Content: "特别好"}, {Content: "顶顶顶"}}},
				{Title: "作者1的第二篇文章", Content: "春晓", Comments: []Comment{{Content: "好诗好诗"}, {Content: "溜溜溜"}}},
			}},
		{Name: "作者2",
			Posts: []Post{
				{Title: "作者2的第一篇文章", Content: "咏鹅", Comments: []Comment{{Content: "生动形象"}, {Content: "传神"}}},
				{Title: "作者2的第二篇文章", Content: "赠汪伦", Comments: []Comment{{Content: "赠离别"}, {Content: "想念"}}},
			}},
		{Name: "作者3",
			Posts: []Post{
				{Title: "作者3的第一篇文章", Content: "咏柳", Comments: []Comment{{Content: "漂亮！"}}},
			}},
	}*/

	comments := []Comment{{Content: "爱你！", PostID: 1},
		{Content: "爱你！", PostID: 1},
		{Content: "爱你！", PostID: 1},
		{Content: "爱你！", PostID: 2},
		{Content: "爱你！", PostID: 2},
		{Content: "爱你！", PostID: 3},
		{Content: "爱你！", PostID: 4},
		{Content: "爱你！", PostID: 3},
		{Content: "爱你！", PostID: 5}}

	tx := repository.DB.Begin()
	if err := tx.Error; err != nil {
		fmt.Println("开启事务失败：", err)
	}

	if err := tx.Create(&comments).Error; err != nil {
		tx.Rollback()
		fmt.Println("新增数据失败：", err)
	}

	tx.Commit()
	fmt.Println("新增数据成功！")

}
