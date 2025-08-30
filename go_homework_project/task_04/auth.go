package task_04

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var JWT_SECRET string

// 用户注册
func Register(c *gin.Context) {
	var regRequest RegisterRequest
	if err := c.ShouldBindJSON(&regRequest); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的请求数据：" + err.Error(),
		})
		return
	}

	//校验用户名是否已存在
	var existUser User
	if err := BlogDB.Where("username =?", regRequest.Username).First(&existUser).Error; err == nil {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "用户名已存在",
		})
		return
	}

	//校验邮箱是否已存在
	if err := BlogDB.Where("email =?", regRequest.Email).First(&existUser).Error; err == nil {
		c.JSON(http.StatusConflict, APIResponse{
			Success: false,
			Error:   "邮箱已存在",
		})
		return
	}

	//密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "密码加密失败",
		})
		return
	}

	//创建用户
	user := User{
		Username: regRequest.Username,
		Password: string(hashedPassword),
		Email:    regRequest.Email,
	}
	if err := BlogDB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "创建用户失败",
		})
	}
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "用户注册成功",
		Data: gin.H{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})

}

// 用户登录
func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "无效的请求数据：" + err.Error(),
		})
		return
	}

	//查询用户信息
	var user User
	if err := BlogDB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "用户名无效",
		})
		return
	}

	//校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Error:   "密码错误",
		})
		return
	}

	//获取token
	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "生成token错误",
		})
		return
	}

	//登陆成功返回信息
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "登陆成功",
		Data: gin.H{
			"token":    token,
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})

}

// 生成token
func generateToken(userID uint, username string) (string, error) {
	//随机生成密钥
	randombytes := make([]byte, 64)
	_, err := rand.Read(randombytes)
	if err != nil {
		Logger.Println("随机生成密钥错误：" + err.Error())
	}

	// 生成一个128字符长的密钥
	JWT_SECRET = hex.EncodeToString(randombytes)
	Logger.Println("JWT_SECRET密钥：", JWT_SECRET)

	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(8 * time.Hour).Unix(), //有效时间8小时
		"iat":      time.Now().Unix(),                    //创建时间
	})

	return token.SignedString([]byte(JWT_SECRET))
}

// JWT鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Authorization header is required",
			})
			c.Abort() //停止当前请求处理链的进一步执行
			return
		}

		//移除前缀 "Bearer ",带空格
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		//解析token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		//获取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid token claims",
			})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid user ID in token",
			})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Error:   "Invalid username in token",
			})
			c.Abort()
			return
		}

		//将用户信息放入上下文
		c.Set("user_id", uint(userID))
		c.Set("username", username)
		c.Next()

	}
}

// 从上下文获取userID
func GetUserId(c *gin.Context) uint {
	userID, exist := c.Get("user_id")
	if !exist {
		return 0
	}
	return userID.(uint)
}

// 从上下文获取username
func GetUserName(c *gin.Context) string {
	username, exist := c.Get("username")
	if !exist {
		return ""
	}
	return username.(string)
}
