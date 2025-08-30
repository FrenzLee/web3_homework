# 个人博客后端系统

这是一个使用Go语言，结合了Gin框架和GORM库开发的一个个人博客系统的后端。实现了博客文章的基本管理功能，包括文章的创建、读取、更新和删除（CRUD）操作，同时支持用户认证和简单的评论功能。

## 主要功能

- 用户的注册和登录（JWT认证）
- 文章的创建、读取、更新、删除（CRUD）
- 评论的创建、读取
- 权限控制（只有作者可以更新、删除自己的文章）
- 文章、评论的分页查询
- 统一的错误处理、响应格式

## 技术栈
- **语言**：Go 1.24.6
- **Web框架**：Gin
- **ORM**：GORM
- **数据库**：MySql
- **认证功能**：JWT
- **加密功能**：bcrypt

## 配置和运行

### 创建数据库
需要在MySql中创建数据库：
```sql
CREATE DATABASE go_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

### 安装依赖

```sybase
go mod tidy
```

### 运行项目
```sybase
go run .
```
项目会在`http://localhost:8080`启动

## API文档

### 基础信息

- **Base URL**: `http://localhost:8080/blogapi`
- **认证方式**: Bearer Token (JWT)
- **响应格式**: JSON

### 认证相关

#### 用户注册

```http
POST /blogapi/auth/register
Content-Type: application/json

{
    "username": "Baby",
    "password": "123456",
    "email" : "Baby@sina.com"
}
```
**响应示例**:
```json
{
  "success": true,
  "message": "用户注册成功",
  "data": {
    "user_id": 1,
    "username": "Baby",
    "email": "Baby@sina.com"
  }
}
```

#### 用户登录

```http
POST /blogapi/auth/login
Content-Type: application/json

{
  "username": "Baby",
  "password": "123456"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "登陆成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 1,
    "username": "Baby",
    "email": "Baby@sina.com"
  }
}
```

### 文章管理

#### 获取文章列表

```http
GET /blogapi/posts?page=1&limit=10
```

#### 获取单个文章

```http
GET /blogapi/posts/{id}
```

#### 创建文章（需要认证）

```http
POST /blogapi/posts
Authorization: Bearer <登录接口返回的token>
Content-Type: application/json

{
  "title": "春晓",
  "content": "春眠不觉晓..."
}
```

#### 更新文章（需要认证）

```http
PUT /blogapi/posts/{id}
Authorization: Bearer <登录接口返回的token>
Content-Type: application/json

{
  "title": "更新的标题",
  "content": "更新的内容..."
}
```


#### 删除文章（需要认证）

```http
DELETE /blogapi/posts/{id}
Authorization: Bearer <登录接口返回的token>
```


### 评论管理

#### 获取文章评论

```http
GET /blogapi/comments/post/{postId}?page=1&limit=20
```

#### 创建评论（需要认证）

```http
POST /blogapi/comments
Authorization: Bearer <登录接口返回的token>
Content-Type: application/json

{
  "content": "评论评论评论",
  "post_id": 1
}
```
## 测试用例

### 使用 Postman 测试

1. **注册用户**
    - Method: POST
    - URL: `http://localhost:8080/blogapi/auth/register`
    - Body:
      ```json
      {
        "username": "Baby",
        "password": "123456",
        "email": "Baby@sina.com"
      }
      ```
2. **登录**
    - Method: POST
    - URL: `http://localhost:8080/blogapi/auth/login`
    - Body:
      ```json
      {
        "username": "Baby",
        "password": "123456"
      }
      ```

3. **创建文章**
    - Method: POST
    - URL: `http://localhost:8080/blogapi/posts`
    - Headers: `Authorization: Bearer <登录接口返回的token>`
    - Body:
      ```json
      {
        "title": "春晓",
        "content": "春眠不觉晓"
      }
      ```
      
4. **查询文章列表**
    - Method: GET
    - URL: `http://localhost:8080/blogapi/posts`

5. **创建评论**
    - Method: POST
    - URL: `http://localhost:8080/blogapi/comments`
    - Headers: `Authorization: Bearer <登录接口返回的token>`
    - Body:
      ```json
      {
        "content": "评论评论",
        "post_id": 1
      }
      ```
      
## 错误处理

系统使用统一的错误响应格式：

```json
{
  "success": true,
  "error": "成功描述信息"
}
```
```json
{
  "success": false,
  "error": "错误描述信息"
}
```

涉及的HTTP状态码：
- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未认证
- `403`: 权限不足
- `404`: 资源不存在
- `409`: 资源冲突（如用户名已存在）
- `500`: 服务器内部错误

## 数据库

项目使用MySQL数据库，需要预先创建数据库。数据库包含以下表：

- `users`: 用户表
- `posts`: 文章表
- `comments`: 评论表

### 数据库配置

默认数据库连接配置：见`config.yaml`文件

## 开发说明

所有需要初始化的相关功能，都在`init.go`中

