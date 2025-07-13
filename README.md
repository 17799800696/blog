# 个人博客系统后端

基于 Go + Gin + GORM 开发的个人博客系统后端API，支持用户认证、文章管理、评论功能等。

## 🚀 功能特性

- ✅ **用户认证与授权** - JWT认证，用户注册登录
- ✅ **文章管理** - 文章的CRUD操作，支持分页
- ✅ **评论系统** - 文章评论功能
- ✅ **权限控制** - 只有作者才能编辑/删除自己的文章
- ✅ **数据库设计** - 完整的数据库模型和关联关系
- ✅ **错误处理** - 统一的错误处理和日志记录
- ✅ **配置管理** - 环境变量配置，支持开发/生产环境
- ✅ **优雅关闭** - 支持优雅关闭和资源清理
- ✅ **请求追踪** - 请求ID追踪和结构化日志
- ✅ **连接池优化** - 数据库连接池配置

## 🛠️ 技术栈

- **语言**: Go 1.24+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT
- **日志**: Zap
- **密码加密**: bcrypt
- **UUID**: Google UUID

## 📋 运行环境要求

- **Go**: 1.24 或更高版本
- **MySQL**: 5.7 或更高版本
- **操作系统**: Linux/macOS/Windows

## 🗄️ 数据库设计

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

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd blog
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置数据库
确保MySQL服务正在运行，并创建数据库：
```sql
CREATE DATABASE blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 启动项目
```bash
./start.sh
```

服务器将在 `http://localhost:8080` 启动

### 5. 测试API
```bash
./test_api.sh
```

## ⚙️ 环境配置

### 必需的环境变量

**服务器配置:**
- `SERVER_PORT`: 服务器端口 (默认: 8080)
- `GIN_MODE`: Gin模式 (debug/release/test)

**数据库配置:**
- `DB_HOST`: 数据库主机 (默认: localhost)
- `DB_PORT`: 数据库端口 (默认: 3306)
- `DB_USER`: 数据库用户名 (默认: root)
- `DB_PASSWORD`: 数据库密码 (默认: 空)
- `DB_NAME`: 数据库名称 (默认: blog)
- `DB_MAX_IDLE_CONNS`: 最大空闲连接数 (默认: 10)
- `DB_MAX_OPEN_CONNS`: 最大打开连接数 (默认: 100)
- `DB_CONN_MAX_LIFETIME`: 连接最大生命周期(分钟) (默认: 60)

**JWT配置:**
- `JWT_SECRET`: JWT密钥 (必需，生产环境必须修改)
- `JWT_EXPIRATION_HOURS`: JWT过期时间（小时）(默认: 24)

**日志配置:**
- `LOG_LEVEL`: 日志级别 (debug/info/warn/error) (默认: info)
- `LOG_FORMAT`: 日志格式 (json/console) (默认: json)
- `LOG_OUTPUT_PATH`: 日志输出路径 (默认: 控制台)

### 开发环境配置
```bash
export SERVER_PORT=8080
export GIN_MODE=debug
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=blog
export DB_MAX_IDLE_CONNS=10
export DB_MAX_OPEN_CONNS=100
export DB_CONN_MAX_LIFETIME=60
export JWT_SECRET=your-secret-key
export JWT_EXPIRATION_HOURS=24
export LOG_LEVEL=info
export LOG_FORMAT=json
export LOG_OUTPUT_PATH=""
```

### 生产环境配置
```bash
export SERVER_PORT=80
export GIN_MODE=release
export DB_HOST=prod-db.example.com
export DB_PORT=3306
export DB_USER=blog_user
export DB_PASSWORD=strong-password
export DB_NAME=blog_prod
export DB_MAX_IDLE_CONNS=20
export DB_MAX_OPEN_CONNS=200
export DB_CONN_MAX_LIFETIME=120
export JWT_SECRET=your-super-secret-jwt-key
export JWT_EXPIRATION_HOURS=168
export LOG_LEVEL=info
export LOG_FORMAT=json
export LOG_OUTPUT_PATH="/var/log/blog/app.log"
```

## 📚 API文档

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

### 文章接口

#### 创建文章 (需要认证)
```http
POST /api/posts
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "title": "文章标题",
  "content": "文章内容"
}
```

#### 获取文章列表
```http
GET /api/posts?page=1&limit=10
```

#### 获取单个文章
```http
GET /api/posts/:id
```

#### 更新文章 (需要认证，仅作者)
```http
PUT /api/posts/:id
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容"
}
```

#### 删除文章 (需要认证，仅作者)
```http
DELETE /api/posts/:id
Authorization: Bearer <your-jwt-token>
```

### 评论接口

#### 创建评论 (需要认证)
```http
POST /api/posts/:id/comments
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "content": "评论内容"
}
```

#### 获取文章评论列表
```http
GET /api/posts/:id/comments?page=1&limit=10
```

### 健康检查
```http
GET /health
```

## 📁 项目结构

```
blog/
├── main.go                    # 主程序入口
├── start.sh                   # 启动脚本
├── test_api.sh               # API测试脚本
├── go.mod                     # Go模块文件
├── go.sum                     # 依赖校验文件
├── README.md                  # 项目说明
├── config/                    # 配置相关
│   ├── config.go             # 配置结构定义
│   ├── database.go           # 数据库配置
│   ├── validator.go          # 配置验证
│   └── README.md             # 配置说明
├── models/                    # 数据模型
│   └── models.go             # 数据库模型定义
├── handlers/                  # 处理器
│   ├── auth.go               # 认证请求结构
│   ├── auth_handler.go       # 认证处理器
│   ├── post.go               # 文章请求结构
│   ├── post_handler.go       # 文章处理器
│   ├── comment.go            # 评论请求结构
│   └── comment_handler.go    # 评论处理器
├── middleware/                # 中间件
│   ├── auth.go               # JWT认证中间件
│   └── request_id.go         # 请求ID中间件
├── routes/                    # 路由配置
│   └── routes.go             # 路由设置
└── utils/                     # 工具函数
    ├── auth.go               # 认证工具
    ├── common.go             # 通用工具
    ├── errors.go             # 错误处理
    └── logger.go             # 日志配置
```

## 🔧 开发指南

### 添加新的API接口

1. **定义请求/响应结构体** (handlers/)
2. **实现处理器函数** (handlers/)
3. **添加路由配置** (routes/routes.go)
4. **编写测试用例** (test_api.sh)

### 数据库迁移

项目使用GORM自动迁移，启动时会自动创建表结构。

### 日志配置

项目使用Zap日志库，支持结构化日志：
- 开发环境：JSON格式，包含时间戳和调用位置
- 生产环境：建议设置LOG_LEVEL=info
- 支持请求ID追踪：每个请求都有唯一的请求ID
- 支持文件输出：通过LOG_OUTPUT_PATH配置

### 优雅关闭

项目支持优雅关闭：
- 监听SIGINT和SIGTERM信号
- 30秒超时关闭
- 自动关闭数据库连接
- 记录关闭日志

## 🧪 测试

### 运行测试脚本
```bash
./test_api.sh
```

测试脚本会依次测试：
1. 健康检查
2. 用户注册
3. 用户登录
4. 获取用户信息
5. 创建文章
6. 获取文章列表
7. 获取单个文章
8. 创建评论
9. 获取评论列表
10. 更新文章
11. 无效token测试
12. 获取不存在的文章
13. 无效注册数据测试
14. 删除文章
15. 文章列表分页
16. 评论列表分页
17. 删除不存在的文章
18. 更新不存在的文章
19. 评论内容验证
20. 文章内容验证

### 手动测试

使用curl或其他HTTP客户端工具测试API：

```bash
# 注册用户
curl -X POST "http://localhost:8080/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456","email":"test@example.com"}'

# 登录获取token
curl -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

## 🔒 安全特性

- ✅ 密码使用bcrypt加密存储
- ✅ JWT token认证
- ✅ 输入验证和错误处理
- ✅ 软删除支持
- ✅ 权限控制（作者只能操作自己的内容）
- ✅ 环境变量配置（敏感信息不硬编码）
- ✅ 请求ID追踪
- ✅ 结构化日志记录

## 🚀 部署

### Docker部署（推荐）

创建Dockerfile：
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### 传统部署

1. 编译项目：`go build -o blog .`
2. 设置环境变量
3. 运行：`./blog`
