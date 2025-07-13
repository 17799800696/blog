#!/bin/bash

# API测试脚本
BASE_URL="http://localhost:8080/api"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试函数
test_api() {
    local test_name="$1"
    local expected_status="$2"
    local response="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # 根据期望状态判断测试结果
    case $expected_status in
        "200"|"201")
            if [[ "$response" == *"\"success\":true"* ]] || [[ "$response" == *"\"status\":\"ok\""* ]]; then
                echo -e "${GREEN}✅ PASS${NC}: $test_name"
                PASSED_TESTS=$((PASSED_TESTS + 1))
            else
                echo -e "${RED}❌ FAIL${NC}: $test_name"
                echo -e "${YELLOW}Response:${NC} $response"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
            ;;
        "400"|"401"|"404")
            if [[ "$response" == *"\"success\":false"* ]] || [[ "$response" == *"\"message\":"* ]]; then
                echo -e "${GREEN}✅ PASS${NC}: $test_name"
                PASSED_TESTS=$((PASSED_TESTS + 1))
            else
                echo -e "${RED}❌ FAIL${NC}: $test_name"
                echo -e "${YELLOW}Response:${NC} $response"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
            ;;
        *)
            # 默认判断成功响应
            if [[ "$response" == *"\"success\":true"* ]] || [[ "$response" == *"\"status\":\"ok\""* ]]; then
                echo -e "${GREEN}✅ PASS${NC}: $test_name"
                PASSED_TESTS=$((PASSED_TESTS + 1))
            else
                echo -e "${RED}❌ FAIL${NC}: $test_name"
                echo -e "${YELLOW}Response:${NC} $response"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
            ;;
    esac
    echo
}

echo -e "${BLUE}=== 博客API测试 ===${NC}"
echo

# 测试健康检查
echo -e "${YELLOW}1. 测试健康检查...${NC}"
HEALTH_RESPONSE=$(curl -s -X GET "$BASE_URL/../health" -H "Content-Type: application/json")
echo "$HEALTH_RESPONSE"
test_api "健康检查" "200" "$HEALTH_RESPONSE"

# 测试用户注册
echo -e "${YELLOW}2. 测试用户注册...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser5",
    "password": "123456",
    "email": "test5@example.com"
  }')
echo "$REGISTER_RESPONSE"
test_api "用户注册" "201" "$REGISTER_RESPONSE"

# 提取token - 修复提取逻辑
TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"//;s/"//')
if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}尝试从登录响应中获取token...${NC}"
    # 如果注册失败，尝试登录获取token
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
      -H "Content-Type: application/json" \
      -d '{
        "username": "testuser",
        "password": "123456"
      }')
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"//;s/"//')
fi

if [ -n "$TOKEN" ]; then
    echo -e "${GREEN}Token获取成功: ${TOKEN:0:20}...${NC}"
else
    echo -e "${RED}Token获取失败${NC}"
fi
echo

# 测试用户登录
echo -e "${YELLOW}3. 测试用户登录...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }')
echo "$LOGIN_RESPONSE"
test_api "用户登录" "200" "$LOGIN_RESPONSE"

# 如果之前没有获取到token，现在重新获取
if [ -z "$TOKEN" ]; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"//;s/"//')
fi

# 测试获取用户信息
echo -e "${YELLOW}4. 测试获取用户信息...${NC}"
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")
echo "$PROFILE_RESPONSE"
test_api "获取用户信息" "200" "$PROFILE_RESPONSE"

# 测试创建文章
echo -e "${YELLOW}5. 测试创建文章...${NC}"
CREATE_POST_RESPONSE=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试文章标题",
    "content": "这是一篇测试文章的内容，用于验证文章创建功能。"
  }')
echo "$CREATE_POST_RESPONSE"
test_api "创建文章" "201" "$CREATE_POST_RESPONSE"

# 提取文章ID
POST_ID=$(echo "$CREATE_POST_RESPONSE" | grep -o '"id":[0-9]*' | sed 's/"id"://')
if [ -n "$POST_ID" ]; then
    echo -e "${GREEN}文章ID: $POST_ID${NC}"
else
    echo -e "${RED}文章ID提取失败${NC}"
    # 使用默认ID进行后续测试
    POST_ID=1
fi
echo

# 测试获取文章列表
echo -e "${YELLOW}6. 测试获取文章列表...${NC}"
POSTS_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/posts" -H "Content-Type: application/json")
echo "$POSTS_LIST_RESPONSE"
test_api "获取文章列表" "200" "$POSTS_LIST_RESPONSE"

# 测试获取单个文章
echo -e "${YELLOW}7. 测试获取单个文章...${NC}"
SINGLE_POST_RESPONSE=$(curl -s -X GET "$BASE_URL/posts/$POST_ID" -H "Content-Type: application/json")
echo "$SINGLE_POST_RESPONSE"
test_api "获取单个文章" "200" "$SINGLE_POST_RESPONSE"

# 测试创建评论
echo -e "${YELLOW}8. 测试创建评论...${NC}"
CREATE_COMMENT_RESPONSE=$(curl -s -X POST "$BASE_URL/posts/$POST_ID/comments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这是一条测试评论，用于验证评论创建功能。"
  }')
echo "$CREATE_COMMENT_RESPONSE"
test_api "创建评论" "201" "$CREATE_COMMENT_RESPONSE"

# 测试获取评论列表
echo -e "${YELLOW}9. 测试获取评论列表...${NC}"
COMMENTS_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/posts/$POST_ID/comments" -H "Content-Type: application/json")
echo "$COMMENTS_LIST_RESPONSE"
test_api "获取评论列表" "200" "$COMMENTS_LIST_RESPONSE"

# 测试更新文章
echo -e "${YELLOW}10. 测试更新文章...${NC}"
UPDATE_POST_RESPONSE=$(curl -s -X PUT "$BASE_URL/posts/$POST_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "更新后的文章标题",
    "content": "这是更新后的文章内容，用于验证文章更新功能。"
  }')
echo "$UPDATE_POST_RESPONSE"
test_api "更新文章" "200" "$UPDATE_POST_RESPONSE"

# 测试错误情况：无效的token
echo -e "${YELLOW}11. 测试无效token...${NC}"
INVALID_TOKEN_RESPONSE=$(curl -s -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer invalid-token" \
  -H "Content-Type: application/json")
echo "$INVALID_TOKEN_RESPONSE"
test_api "无效token测试" "401" "$INVALID_TOKEN_RESPONSE"

# 测试错误情况：不存在的文章
echo -e "${YELLOW}12. 测试获取不存在的文章...${NC}"
NOT_FOUND_RESPONSE=$(curl -s -X GET "$BASE_URL/posts/99999" -H "Content-Type: application/json")
echo "$NOT_FOUND_RESPONSE"
test_api "获取不存在的文章" "404" "$NOT_FOUND_RESPONSE"

# 测试错误情况：无效的请求数据
echo -e "${YELLOW}13. 测试无效的注册数据...${NC}"
INVALID_REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "password": "123",
    "email": "invalid-email"
  }')
echo "$INVALID_REGISTER_RESPONSE"
test_api "无效注册数据测试" "400" "$INVALID_REGISTER_RESPONSE"

# 测试删除文章
echo -e "${YELLOW}14. 测试删除文章...${NC}"
DELETE_POST_RESPONSE=$(curl -s -X DELETE "$BASE_URL/posts/$POST_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")
echo "$DELETE_POST_RESPONSE"
test_api "删除文章" "200" "$DELETE_POST_RESPONSE"

# 测试分页功能
echo -e "${YELLOW}15. 测试文章列表分页...${NC}"
PAGINATION_RESPONSE=$(curl -s -X GET "$BASE_URL/posts?page=1&limit=5" -H "Content-Type: application/json")
echo "$PAGINATION_RESPONSE"
test_api "文章列表分页" "200" "$PAGINATION_RESPONSE"

# 测试评论分页功能 - 使用存在的文章ID
echo -e "${YELLOW}16. 测试评论列表分页...${NC}"
COMMENT_PAGINATION_RESPONSE=$(curl -s -X GET "$BASE_URL/posts/5/comments?page=1&limit=5" -H "Content-Type: application/json")
echo "$COMMENT_PAGINATION_RESPONSE"
test_api "评论列表分页" "200" "$COMMENT_PAGINATION_RESPONSE"

# 测试权限控制：尝试删除不存在的文章
echo -e "${YELLOW}17. 测试删除不存在的文章...${NC}"
DELETE_NOT_FOUND_RESPONSE=$(curl -s -X DELETE "$BASE_URL/posts/99999" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")
echo "$DELETE_NOT_FOUND_RESPONSE"
test_api "删除不存在的文章" "404" "$DELETE_NOT_FOUND_RESPONSE"

# 测试权限控制：尝试更新不存在的文章
echo -e "${YELLOW}18. 测试更新不存在的文章...${NC}"
UPDATE_NOT_FOUND_RESPONSE=$(curl -s -X PUT "$BASE_URL/posts/99999" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "更新不存在的文章",
    "content": "这应该失败"
  }')
echo "$UPDATE_NOT_FOUND_RESPONSE"
test_api "更新不存在的文章" "404" "$UPDATE_NOT_FOUND_RESPONSE"

# 测试评论内容验证
echo -e "${YELLOW}19. 测试评论内容验证...${NC}"
INVALID_COMMENT_RESPONSE=$(curl -s -X POST "$BASE_URL/posts/1/comments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": ""
  }')
echo "$INVALID_COMMENT_RESPONSE"
test_api "评论内容验证" "400" "$INVALID_COMMENT_RESPONSE"

# 测试文章内容验证
echo -e "${YELLOW}20. 测试文章内容验证...${NC}"
INVALID_POST_RESPONSE=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "",
    "content": ""
  }')
echo "$INVALID_POST_RESPONSE"
test_api "文章内容验证" "400" "$INVALID_POST_RESPONSE"

# 输出测试结果统计
echo -e "${BLUE}=== 测试结果统计 ===${NC}"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"
echo -e "${BLUE}总计: $TOTAL_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
else
    echo -e "${RED}❌ 有 $FAILED_TESTS 个测试失败${NC}"
fi

echo -e "${BLUE}=== 测试完成 ===${NC}" 