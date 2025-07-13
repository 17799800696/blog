#!/bin/bash

# 博客系统启动脚本

echo "Setting up environment variables..."

# 服务器配置
export SERVER_PORT=8080
export GIN_MODE=debug

# 数据库配置
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=123456
export DB_NAME=blog

# 数据库连接池配置
export DB_MAX_IDLE_CONNS=10
export DB_MAX_OPEN_CONNS=100
export DB_CONN_MAX_LIFETIME=60

# JWT配置
export JWT_SECRET=your-super-secret-jwt-key-change-in-production
export JWT_EXPIRATION_HOURS=24

# 日志配置
export LOG_LEVEL=info
export LOG_FORMAT=json
export LOG_OUTPUT_PATH=""

echo "Environment variables set successfully!"
echo "Starting blog server..."

# 启动服务器
go run main.go 