#!/bin/bash

# API 测试脚本
# 使用方法: ./test_api.sh

BASE_URL="http://localhost:8080"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}   Go ViewSet API 测试脚本${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 健康检查
echo -e "${GREEN}=== 1. 健康检查 ===${NC}"
curl -s $BASE_URL/health | jq .
echo -e "\n"

# 创建用户
echo -e "${GREEN}=== 2. 创建用户 ===${NC}"
curl -s -X POST $BASE_URL/api/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试用户",
    "email": "test@example.com",
    "status": "active",
    "age": 22,
    "phone": "13900139000"
  }' | jq .
echo -e "\n"

# 获取用户列表
echo -e "${GREEN}=== 3. 获取用户列表（分页） ===${NC}"
curl -s "$BASE_URL/api/users/?page=1&page_size=5" | jq .
echo -e "\n"

# 过滤查询
echo -e "${GREEN}=== 4. 过滤查询（status=active） ===${NC}"
curl -s "$BASE_URL/api/users/?status=active" | jq .
echo -e "\n"

# 排序查询
echo -e "${GREEN}=== 5. 排序查询（按年龄降序） ===${NC}"
curl -s "$BASE_URL/api/users/?order_by=age desc" | jq .
echo -e "\n"

# 获取单个用户
echo -e "${GREEN}=== 6. 获取单个用户 ===${NC}"
curl -s $BASE_URL/api/users/1 | jq .
echo -e "\n"

# 更新用户
echo -e "${GREEN}=== 7. 更新用户 ===${NC}"
curl -s -X PUT $BASE_URL/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三（已修改）",
    "age": 27
  }' | jq .
echo -e "\n"

# 激活用户
echo -e "${GREEN}=== 8. 激活用户（自定义 Action） ===${NC}"
curl -s -X POST $BASE_URL/api/users/2/activate | jq .
echo -e "\n"

# 停用用户
echo -e "${GREEN}=== 9. 停用用户（自定义 Action） ===${NC}"
curl -s -X POST $BASE_URL/api/users/1/deactivate | jq .
echo -e "\n"

# 重置密码
echo -e "${GREEN}=== 10. 重置密码（自定义 Action） ===${NC}"
curl -s -X POST $BASE_URL/api/users/1/reset_password | jq .
echo -e "\n"

# 获取统计信息
echo -e "${GREEN}=== 11. 获取统计信息 ===${NC}"
curl -s $BASE_URL/api/users/stats | jq .
echo -e "\n"

# 删除用户
echo -e "${GREEN}=== 12. 删除用户 ===${NC}"
curl -s -X DELETE $BASE_URL/api/users/3 | jq .
echo -e "\n"

# 错误测试：获取不存在的用户
echo -e "${YELLOW}=== 13. 错误测试：获取不存在的用户 ===${NC}"
curl -s $BASE_URL/api/users/999 | jq .
echo -e "\n"

# 错误测试：创建重复邮箱
echo -e "${YELLOW}=== 14. 错误测试：创建重复邮箱 ===${NC}"
curl -s -X POST $BASE_URL/api/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "重复用户",
    "email": "zhangsan@example.com"
  }' | jq .
echo -e "\n"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   测试完成！${NC}"
echo -e "${GREEN}========================================${NC}"
