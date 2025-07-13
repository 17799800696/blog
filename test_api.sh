#!/bin/bash

# API测试脚本
BASE_URL="http://localhost:8080/api"

echo "=== 博客API测试 ==="

# 测试健康检查
echo "1. 测试健康检查..."
curl -X GET "$BASE_URL/../health" -H "Content-Type: application/json"
echo -e "\n"

# 测试用户注册
echo "2. 测试用户注册..."
curl -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com"
  }'
echo -e "\n"

# 测试用户登录
echo "3. 测试用户登录..."
curl -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
echo -e "\n"

echo "=== 测试完成 ===" 