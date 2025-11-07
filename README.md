# Go ViewSet - Django Rest Framework 风格的 Go 封装

基于 Gin + GORM 的 RESTful API 快速开发框架，灵感来自 Django Rest Framework。

## 特性

- 🚀 快速生成 RESTful CRUD 接口
- 📦 模块化架构：Controller / Service / Model / Router
- 🔍 自动分页、过滤、排序
- 🎯 自定义 Action 支持（类似 DRF 的 @action）
- 📝 统一的响应格式
- 🛡️ 统一的错误处理
- 🔌 易于扩展和定制

## 快速开始

### 安装依赖

```bash
go mod download
```

### 运行示例

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 示例

### 1. 创建用户
```bash
curl -X POST http://localhost:8080/api/users/ \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com","status":"active"}'
```

### 2. 获取用户列表（支持分页和过滤）
```bash
# 基础列表
curl http://localhost:8080/api/users/

# 分页
curl "http://localhost:8080/api/users/?page=1&page_size=10"

# 过滤
curl "http://localhost:8080/api/users/?status=active&name=张三"

# 排序
curl "http://localhost:8080/api/users/?order_by=created_at desc"
```

### 3. 获取单个用户
```bash
curl http://localhost:8080/api/users/1
```

### 4. 更新用户
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"李四","email":"lisi@example.com"}'
```

### 5. 删除用户
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

### 6. 自定义 Action
```bash
# 激活用户
curl -X POST http://localhost:8080/api/users/1/activate

# 重置密码
curl -X POST http://localhost:8080/api/users/1/reset_password
```

## 项目结构

```
go-viewset/
├── go.mod                          # Go 模块依赖
├── main.go                         # 主程序入口
└── internal/
    ├── models/                     # 数据模型
    │   └── user.go
    ├── viewset/                    # ViewSet 层
    │   ├── base_viewset.go        # 基础 ViewSet
    │   └── user_viewset.go        # 用户 ViewSet
    ├── utils/                      # 工具函数
    │   ├── response.go            # 统一响应格式
    │   ├── pagination.go          # 分页工具
    │   └── filter.go              # 过滤和排序工具
    └── router/
        └── router.go              # 路由注册
```

## 核心概念

### ViewSet

ViewSet 是一个封装了标准 CRUD 操作的控制器。通过嵌入 `GenericViewSet`，你可以快速创建 RESTful API。

```go
type UserViewSet struct {
    *viewset.GenericViewSet
}

func NewUserViewSet(db *gorm.DB) *UserViewSet {
    return &UserViewSet{
        GenericViewSet: viewset.NewGenericViewSet(db, &models.User{}),
    }
}
```

### 自定义 Action

使用 `RegisterAction` 方法可以注册自定义操作：

```go
func (v *UserViewSet) RegisterRoutes(group *gin.RouterGroup) {
    v.GenericViewSet.RegisterRoutes(group)
    
    // 注册自定义 action
    v.RegisterAction("POST", "/:id/activate", v.Activate)
}

func (v *UserViewSet) Activate(c *gin.Context) {
    // 自定义逻辑
}
```

### 过滤和排序

框架自动解析查询参数：

- `?name=value` - 等值过滤
- `?order_by=field desc` - 排序
- `?page=1&page_size=10` - 分页

### 统一响应格式

所有接口返回统一的 JSON 格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {...},
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 100
  }
}
```

## 扩展你的 ViewSet

### 1. 创建模型

```go
type Product struct {
    gorm.Model
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Stock int     `json:"stock"`
}
```

### 2. 创建 ViewSet

```go
type ProductViewSet struct {
    *viewset.GenericViewSet
}

func NewProductViewSet(db *gorm.DB) *ProductViewSet {
    return &ProductViewSet{
        GenericViewSet: viewset.NewGenericViewSet(db, &Product{}),
    }
}
```

### 3. 注册路由

```go
productViewSet := NewProductViewSet(db)
productViewSet.RegisterRoutes(r.Group("/api/products"))
```

就这么简单！

## 进阶功能

### 重写默认方法

你可以重写任何默认方法来自定义行为：

```go
func (v *ProductViewSet) Create(c *gin.Context) {
    // 自定义创建逻辑
    // 例如：添加额外验证、发送通知等
}
```

### 添加中间件

```go
productViewSet.RegisterRoutes(
    r.Group("/api/products").Use(AuthMiddleware()),
)
```

### 自定义查询

在 ViewSet 中可以访问 `v.DB` 进行自定义查询：

```go
func (v *UserViewSet) GetActiveUsers(c *gin.Context) {
    var users []models.User
    v.DB.Where("status = ?", "active").Find(&users)
    utils.Success(c, users)
}
```

## 技术栈

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **数据库**: SQLite（示例用，可替换为 MySQL/PostgreSQL）

## License

MIT License
