# 个人博客系统后端

基于 Go + Gin + GORM 开发的个人博客系统后端API。

## 功能特性

- ✅ 用户认证与授权（JWT）
- ✅ 用户注册和登录
- ✅ 密码加密存储（bcrypt）
- ✅ 数据库设计与模型定义
- 🔄 文章CRUD操作（待实现）
- 🔄 评论功能（待实现）

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **密码加密**: bcrypt

## 项目结构

```
blog/
├── main.go              # 主程序入口
├── config/
│   └── database.go      # 数据库配置
├── models/
│   └── models.go        # 数据模型定义
├── handlers/
│   ├── auth.go          # 认证相关请求结构
│   └── auth_handler.go  # 认证处理器
├── middleware/
│   └── auth.go          # JWT认证中间件
├── routes/
│   └── routes.go        # 路由配置
├── utils/
│   └── auth.go          # 认证工具函数
└── test_api.sh          # API测试脚本
```

## 数据库设计

### users 表
- `id` (主键)
- `username` (用户名，唯一)
- `password` (加密密码)
- `email` (邮箱，唯一)
- `created_at`, `updated_at`, `deleted_at`

### posts 表
- `id` (主键)
- `title` (文章标题)
- `content` (文章内容)
- `user_id` (关联用户)
- `created_at`, `updated_at`, `deleted_at`

### comments 表
- `id` (主键)
- `content` (评论内容)
- `user_id` (关联用户)
- `post_id` (关联文章)
- `created_at`, `deleted_at`

## 环境配置

### 数据库配置
通过环境变量配置数据库连接：

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=blog
```

### 默认配置
- 数据库主机: localhost
- 端口: 3306
- 用户名: root
- 密码: 123456
- 数据库名: blog

## 运行项目

1. **安装依赖**
```bash
go mod tidy
```

2. **启动服务器**
```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

3. **测试API**
```bash
./test_api.sh
```

## API接口

### 认证接口

#### 用户注册
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com"
}
```

#### 用户登录
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

#### 获取用户信息
```http
GET /api/profile
Authorization: Bearer <your-jwt-token>
```

### 健康检查
```http
GET /health
```

## 安全特性

- ✅ 密码使用 bcrypt 加密存储
- ✅ JWT token 认证
- ✅ 软删除支持
- ✅ 输入验证和错误处理
- ✅ 唯一性约束（用户名、邮箱）

## 下一步计划

- [ ] 实现文章的CRUD操作
- [ ] 实现评论功能
- [ ] 添加文章分类和标签
- [ ] 实现文件上传功能
- [ ] 添加API文档（Swagger）
- [ ] 添加单元测试 