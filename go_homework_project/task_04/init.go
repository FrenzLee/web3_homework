package task_04

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Logger *log.Logger
var BlogDB *gorm.DB
var Router *gin.Engine

// 初始化Logger
func initLogger() {
	//打开日志文件
	logfile, err := os.OpenFile("blog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个新的logger实例，设置前缀和标志位
	Logger = log.New(logfile, "blog log：", log.LstdFlags)

}

// 读取配置文件
func loadConfig(filename string, configDB *ConfigDB) error {
	//打开文件
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("无法打开配置文件 %d：%v", filename, err)
	}
	defer file.Close()

	// 读取文件内容到字节切片
	var bytes []byte
	bytes, err = os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("读取配置文件失败 %d：%v", filename, err)
	}

	// 解析YAML内容到结构体
	err = yaml.Unmarshal(bytes, configDB)
	if err != nil {
		return fmt.Errorf("解析配置文件失败 %d：%v", filename, err)
	}
	return nil
}

func initDB() {
	var configDB ConfigDB

	//获取当前文件运行目录
	pwd, err := os.Getwd()
	if err != nil {
		Logger.Fatalf("无法获取当前文件目录：%v", err)
		return
	}
	configPath := filepath.Join(pwd, "task_04", "config.yaml")

	//加载配置
	if err := loadConfig(configPath, &configDB); err != nil {
		Logger.Fatalf("无法加载配置：%v", err)
	}

	// 创建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configDB.Database.Username,
		configDB.Database.Password,
		configDB.Database.Host,
		configDB.Database.Port,
		configDB.Database.DBName,
	)
	Logger.Printf("数据库连接字符串：%s", dsn)

	// 打开数据库连接
	BlogDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		Logger.Fatalf("无法连接到数据库：%v", err)
	}

	// 测试连接是否成功
	sqlDB, _ := BlogDB.DB()
	if err := sqlDB.Ping(); err != nil {
		Logger.Fatalf("数据库连接失败：%v", err)
	}

	Logger.Println("数据库连接成功")
}

// 创建表
func initTable() {
	if BlogDB == nil {
		Logger.Fatalf("BlogDB 未初始化")
	}

	if err := BlogDB.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		Logger.Fatalf("创建表失败：%v", err)
	}
	Logger.Println("表创建成功")
}

// 初始化路由
func initRouter() {

	//创建路由
	Router = gin.Default()

	//添加中间件
	Router.Use(gin.Logger())   //为每个请求记录一些基本信息，如请求时间、HTTP方法、请求路径、响应状态码等
	Router.Use(gin.Recovery()) //捕获任何在处理请求过程中发生的 panic，并将其恢复，以防止整个程序崩溃。它还会记录引发 panic 的堆栈信息。

	// 设置路由组
	blog_api := Router.Group("/blogapi")

	//用户
	auth := blog_api.Group("/auth")
	auth.POST("/register", Register)
	auth.POST("/login", Login)

	//文章
	post := blog_api.Group("/posts")
	post.GET("", GetPosts)                            //获取所有文章列表
	post.GET("/:id", GetPost)                         //获取单个文章详情
	post.POST("", AuthMiddleware(), CreatePost)       //创建新文章
	post.PUT("/:id", AuthMiddleware(), UpdatePost)    //更新文章
	post.DELETE("/:id", AuthMiddleware(), DeletePost) //删除文章

	//评论
	comment := blog_api.Group("/comments")
	comment.GET("/:postId", GetCommentsByPost)               //获取指定文章的所有评论
	comment.POST("", AuthMiddleware(), CreateCommentForPost) //创建新评论

}

func Init() {
	//初始化Logger
	initLogger()

	//初始化BlogDB
	initDB()

	//创建表
	//initTable()

	// 初始化路由
	initRouter()
}
