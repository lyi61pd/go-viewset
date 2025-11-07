# 使用示例

本文档提供完整的 API 使用示例，包括 curl 命令和返回结果。

## 启动服务

```bash
cd /Users/ybbj100324/code/go-viewset
go mod download
go run main.go
```

服务将在 `http://localhost:8080` 启动。

---

## API 测试示例

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

**响应:**
```json
{
  "status": "ok",
  "message": "Go ViewSet is running"
}
```

---

### 2. 创建用户 (POST /api/users/)

```bash
curl -X POST http://localhost:8080/api/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "赵六",
    "email": "zhaoliu@example.com",
    "status": "active",
    "age": 26,
    "phone": "13800138003"
  }'
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 4,
    "created_at": "2025-11-07T10:30:00Z",
    "updated_at": "2025-11-07T10:30:00Z",
    "name": "赵六",
    "email": "zhaoliu@example.com",
    "status": "active",
    "age": 26,
    "phone": "13800138003"
  }
}
```

**错误示例（邮箱重复）:**
```bash
curl -X POST http://localhost:8080/api/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "status": "active"
  }'
```

**响应:**
```json
{
  "code": 400,
  "msg": "该邮箱已被注册"
}
```

---

### 3. 获取用户列表 (GET /api/users/)

#### 3.1 基础列表

```bash
curl http://localhost:8080/api/users/
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "created_at": "2025-11-07T10:00:00Z",
      "updated_at": "2025-11-07T10:00:00Z",
      "name": "张三",
      "email": "zhangsan@example.com",
      "status": "active",
      "age": 25,
      "phone": "13800138000"
    },
    {
      "id": 2,
      "created_at": "2025-11-07T10:00:00Z",
      "updated_at": "2025-11-07T10:00:00Z",
      "name": "李四",
      "email": "lisi@example.com",
      "status": "active",
      "age": 30,
      "phone": "13800138001"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 3
  }
}
```

#### 3.2 分页查询

```bash
curl "http://localhost:8080/api/users/?page=1&page_size=2"
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 2,
    "total": 3
  }
}
```

#### 3.3 过滤查询

```bash
# 按状态过滤
curl "http://localhost:8080/api/users/?status=active"

# 按名称过滤
curl "http://localhost:8080/api/users/?name=张三"

# 多条件过滤
curl "http://localhost:8080/api/users/?status=active&age=25"
```

#### 3.4 排序查询

```bash
# 按创建时间降序
curl "http://localhost:8080/api/users/?order_by=created_at desc"

# 按年龄升序
curl "http://localhost:8080/api/users/?order_by=age asc"

# DRF 风格排序（-表示降序）
curl "http://localhost:8080/api/users/?ordering=-created_at"
```

#### 3.5 组合查询

```bash
curl "http://localhost:8080/api/users/?status=active&page=1&page_size=5&order_by=age desc"
```

---

### 4. 获取单个用户 (GET /api/users/:id)

```bash
curl http://localhost:8080/api/users/1
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "created_at": "2025-11-07T10:00:00Z",
    "updated_at": "2025-11-07T10:00:00Z",
    "name": "张三",
    "email": "zhangsan@example.com",
    "status": "active",
    "age": 25,
    "phone": "13800138000"
  }
}
```

**错误示例（用户不存在）:**
```bash
curl http://localhost:8080/api/users/999
```

**响应:**
```json
{
  "code": 404,
  "msg": "记录不存在"
}
```

---

### 5. 更新用户 (PUT /api/users/:id)

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三（已更新）",
    "email": "zhangsan_new@example.com",
    "age": 26
  }'
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "created_at": "2025-11-07T10:00:00Z",
    "updated_at": "2025-11-07T10:35:00Z",
    "name": "张三（已更新）",
    "email": "zhangsan_new@example.com",
    "status": "active",
    "age": 26,
    "phone": "13800138000"
  }
}
```

---

### 6. 删除用户 (DELETE /api/users/:id)

```bash
curl -X DELETE http://localhost:8080/api/users/3
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "删除成功"
  }
}
```

---

### 7. 自定义 Action

#### 7.1 激活用户 (POST /api/users/:id/activate)

```bash
curl -X POST http://localhost:8080/api/users/2/activate
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "用户已激活",
    "user": {
      "id": 2,
      "created_at": "2025-11-07T10:00:00Z",
      "updated_at": "2025-11-07T10:40:00Z",
      "name": "李四",
      "email": "lisi@example.com",
      "status": "active",
      "age": 30,
      "phone": "13800138001"
    }
  }
}
```

#### 7.2 停用用户 (POST /api/users/:id/deactivate)

```bash
curl -X POST http://localhost:8080/api/users/1/deactivate
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "用户已停用",
    "user": {
      "id": 1,
      "status": "inactive",
      ...
    }
  }
}
```

#### 7.3 重置密码 (POST /api/users/:id/reset_password)

```bash
curl -X POST http://localhost:8080/api/users/1/reset_password
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "密码重置邮件已发送",
    "user_id": 1,
    "email": "zhangsan@example.com"
  }
}
```

#### 7.4 获取统计信息 (GET /api/users/stats)

```bash
curl http://localhost:8080/api/users/stats
```

**响应:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 3,
    "active": 2,
    "inactive": 1
  }
}
```

---

## 测试脚本

你可以创建一个测试脚本来批量测试所有 API：

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== 1. 健康检查 ==="
curl $BASE_URL/health
echo -e "\n"

echo "=== 2. 创建用户 ==="
curl -X POST $BASE_URL/api/users/ \
  -H "Content-Type: application/json" \
  -d '{"name":"测试用户","email":"test@example.com","status":"active","age":22}'
echo -e "\n"

echo "=== 3. 获取用户列表 ==="
curl "$BASE_URL/api/users/?page=1&page_size=5"
echo -e "\n"

echo "=== 4. 过滤查询（status=active） ==="
curl "$BASE_URL/api/users/?status=active"
echo -e "\n"

echo "=== 5. 获取单个用户 ==="
curl $BASE_URL/api/users/1
echo -e "\n"

echo "=== 6. 更新用户 ==="
curl -X PUT $BASE_URL/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"张三（已修改）","age":27}'
echo -e "\n"

echo "=== 7. 激活用户 ==="
curl -X POST $BASE_URL/api/users/1/activate
echo -e "\n"

echo "=== 8. 获取统计信息 ==="
curl $BASE_URL/api/users/stats
echo -e "\n"

echo "=== 测试完成 ==="
```

保存为 `test_api.sh`，然后执行：

```bash
chmod +x test_api.sh
./test_api.sh
```

---

## 使用 Postman / Insomnia

你也可以导入以下 API 集合到 Postman 或 Insomnia 中进行测试：

### 请求集合

1. **创建用户**
   - Method: POST
   - URL: `http://localhost:8080/api/users/`
   - Body: `{"name":"测试","email":"test@example.com","status":"active"}`

2. **获取列表**
   - Method: GET
   - URL: `http://localhost:8080/api/users/?page=1&page_size=10`

3. **获取详情**
   - Method: GET
   - URL: `http://localhost:8080/api/users/1`

4. **更新用户**
   - Method: PUT
   - URL: `http://localhost:8080/api/users/1`
   - Body: `{"name":"更新后的名称"}`

5. **删除用户**
   - Method: DELETE
   - URL: `http://localhost:8080/api/users/1`

6. **自定义 Action**
   - Method: POST
   - URL: `http://localhost:8080/api/users/1/activate`

---

## 常见问题

### Q1: 如何清空数据库？

删除 `test.db` 文件，然后重新运行程序：

```bash
rm test.db
go run main.go
```

### Q2: 如何修改端口？

修改 `main.go` 中的 `port` 变量：

```go
port := ":8080"  // 改为你想要的端口
```

### Q3: 如何使用 MySQL 或 PostgreSQL？

修改 `main.go` 中的数据库连接：

```go
// MySQL
import "gorm.io/driver/mysql"
dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// PostgreSQL
import "gorm.io/driver/postgres"
dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```
