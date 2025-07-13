# 配置说明

## 必需的环境变量

**所有环境变量都是必需的，没有默认值！**

### 服务器配置
- `SERVER_PORT`: 服务器端口 (必需)
- `GIN_MODE`: Gin模式 (必需，可选: debug, release, test)

### 数据库配置
- `DB_HOST`: 数据库主机 (必需)
- `DB_PORT`: 数据库端口 (必需)
- `DB_USER`: 数据库用户名 (必需)
- `DB_PASSWORD`: 数据库密码 (必需)
- `DB_NAME`: 数据库名称 (必需)

### JWT配置
- `JWT_SECRET`: JWT密钥 (必需，生产环境必须修改)
- `JWT_EXPIRATION_HOURS`: JWT过期时间（小时）(必需，必须大于0)

## 启动方式

### 方式1：使用启动脚本（推荐）
```bash
./start.sh
```

### 方式2：手动设置环境变量
```bash
# 服务器配置
export SERVER_PORT=8080
export GIN_MODE=debug

# 数据库配置
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=blog

# JWT配置
export JWT_SECRET=your-super-secret-jwt-key-change-in-production
export JWT_EXPIRATION_HOURS=24

# 启动服务器
go run main.go
```

### 方式3：使用.env文件（需要安装godotenv）
```bash
# 安装godotenv
go get github.com/joho/godotenv

# 创建.env文件
cat > .env << EOF
SERVER_PORT=8080
GIN_MODE=debug
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=blog
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION_HOURS=24
EOF

# 启动服务器
go run main.go
```

## 配置验证

启动时会自动验证：
1. 所有必需的环境变量是否已设置
2. JWT过期时间是否大于0
3. Gin模式是否有效
4. 服务器端口是否有效

如果验证失败，程序会立即退出并显示错误信息。

## 生产环境配置

```bash
export SERVER_PORT=80
export GIN_MODE=release
export DB_HOST=prod-db.example.com
export DB_PORT=3306
export DB_USER=blog_user
export DB_PASSWORD=strong-password-here
export DB_NAME=blog_prod
export JWT_SECRET=your-super-secret-jwt-key-change-in-production
export JWT_EXPIRATION_HOURS=168  # 7天
```

## 安全建议

1. **生产环境必须修改 JWT_SECRET**
2. **使用强密码保护数据库**
3. **设置合适的 JWT 过期时间**
4. **使用 HTTPS 和防火墙保护服务器**
5. **定期轮换 JWT 密钥**
6. **使用环境变量而不是硬编码配置** 