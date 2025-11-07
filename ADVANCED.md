# 进阶功能指南

本文档介绍如何扩展和增强 Go ViewSet 框架，添加更多企业级功能。

---

## 1. 配置管理

### 使用 Viper 进行配置管理

**安装依赖:**
```bash
go get github.com/spf13/viper
```

**创建配置文件 `config.yaml`:**
```yaml
server:
  port: 8080
  mode: debug  # debug, release, test

database:
  driver: mysql
  host: localhost
  port: 3306
  database: mydb
  username: root
  password: password
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

jwt:
  secret: your-secret-key
  expire_hours: 24

log:
  level: info
  file: logs/app.log
```

**创建配置管理器 `internal/config/config.go`:**
```go
package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Database     string `mapstructure:"database"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &config, nil
}
```

---

## 2. JWT 认证

### 实现 JWT 中间件

**安装依赖:**
```bash
go get github.com/golang-jwt/jwt/v5
```

**创建 JWT 工具 `internal/utils/jwt.go`:**
```go
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key")

// GenerateToken 生成 JWT token
func GenerateToken(userID uint, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
```

**创建认证中间件 `internal/middleware/auth.go`:**
```go
package middleware

import (
	"go-viewset/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 验证 token 格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户 ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}
```

**使用认证中间件:**
```go
// 在 router.go 中
protectedAPI := api.Group("/")
protectedAPI.Use(middleware.AuthMiddleware())
{
    userViewSet := viewset.NewUserViewSet(db)
    userViewSet.RegisterRoutes(protectedAPI.Group("/users"))
}
```

---

## 3. 权限控制

### 实现基于角色的访问控制 (RBAC)

**创建权限模型 `internal/models/role.go`:**
```go
package models

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string       `gorm:"size:200" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:200" json:"description"`
	Resource    string `gorm:"size:50" json:"resource"` // 资源名称，如 "users"
	Action      string `gorm:"size:20" json:"action"`   // 操作，如 "read", "write", "delete"
}

// 为 User 模型添加角色关联
type User struct {
	// ... 其他字段
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}
```

**创建权限中间件 `internal/middleware/permission.go`:**
```go
package middleware

import (
	"go-viewset/internal/models"
	"go-viewset/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequirePermission 权限检查中间件
func RequirePermission(db *gorm.DB, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetCurrentUserID(c)
		if !exists {
			utils.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 查询用户及其角色权限
		var user models.User
		if err := db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
			utils.InternalServerError(c, "查询用户失败")
			c.Abort()
			return
		}

		// 检查是否有所需权限
		hasPermission := false
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				if perm.Resource == resource && perm.Action == action {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			utils.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}
```

**使用权限中间件:**
```go
// 只有有 "users:write" 权限的用户才能创建用户
api.POST("/users", middleware.RequirePermission(db, "users", "write"), userViewSet.Create)
```

---

## 4. 数据验证

### 使用 validator 进行高级验证

**创建自定义验证器 `internal/validator/custom.go`:**
```go
package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("username", validateUsername)
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	match, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return match
}

// validateUsername 验证用户名（只能包含字母、数字、下划线）
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
	return match
}
```

**使用自定义验证器:**
```go
type User struct {
	// ...
	Phone    string `json:"phone" binding:"phone"`
	Username string `json:"username" binding:"required,username"`
}
```

---

## 5. 日志管理

### 使用 Zap 实现结构化日志

**安装依赖:**
```bash
go get go.uber.org/zap
```

**创建日志工具 `internal/utils/logger.go`:**
```go
package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger 初始化日志
func InitLogger(level string, logFile string) error {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logFile}
	
	// 设置日志级别
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}

	var err error
	Logger, err = config.Build()
	return err
}
```

**在代码中使用:**
```go
utils.Logger.Info("用户创建成功", 
	zap.Uint("user_id", user.ID),
	zap.String("email", user.Email),
)
```

---

## 6. Swagger 文档

### 使用 swaggo 生成 API 文档

**安装依赖:**
```bash
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
go install github.com/swaggo/swag/cmd/swag@latest
```

**添加 Swagger 注释:**
```go
// @title           Go ViewSet API
// @version         1.0
// @description     类似 Django Rest Framework 的 Go 封装
// @host            localhost:8080
// @BasePath        /api

// @Summary      创建用户
// @Description  创建一个新用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "用户信息"
// @Success      200   {object}  utils.Response
// @Router       /users/ [post]
func (v *UserViewSet) Create(c *gin.Context) {
	// ...
}
```

**生成文档:**
```bash
swag init
```

**注册 Swagger 路由:**
```go
import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

访问 `http://localhost:8080/swagger/index.html` 查看文档。

---

## 7. 缓存支持

### 使用 Redis 缓存

**安装依赖:**
```bash
go get github.com/redis/go-redis/v9
```

**创建缓存工具 `internal/utils/cache.go`:**
```go
package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis
func InitRedis(addr, password string, db int) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	return err
}

// CacheGet 获取缓存
func CacheGet(key string, dest interface{}) error {
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// CacheSet 设置缓存
func CacheSet(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.Set(context.Background(), key, data, expiration).Err()
}
```

**在 ViewSet 中使用缓存:**
```go
func (v *UserViewSet) Retrieve(c *gin.Context) {
	id := c.Param("id")
	cacheKey := fmt.Sprintf("user:%s", id)

	// 先尝试从缓存获取
	var user models.User
	if err := utils.CacheGet(cacheKey, &user); err == nil {
		utils.Success(c, user)
		return
	}

	// 缓存未命中，从数据库查询
	if err := v.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 写入缓存
	utils.CacheSet(cacheKey, user, 10*time.Minute)

	utils.Success(c, user)
}
```

---

## 8. 数据库迁移

### 使用 golang-migrate

**安装工具:**
```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**创建迁移文件:**
```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

**编写迁移 SQL:**
```sql
-- migrations/000001_create_users_table.up.sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'inactive',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- migrations/000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

**执行迁移:**
```bash
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/dbname" up
```

---

## 9. 测试

### 编写单元测试

**创建测试文件 `internal/viewset/user_viewset_test.go`:**
```go
package viewset_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})
	return db
}

func TestUserViewSet_Create(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	
	viewSet := NewUserViewSet(db)
	viewSet.RegisterRoutes(router.Group("/users"))

	user := map[string]interface{}{
		"name":  "测试用户",
		"email": "test@example.com",
	}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
```

**运行测试:**
```bash
go test ./...
```

---

## 10. Docker 部署

**创建 Dockerfile:**
```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./main"]
```

**创建 docker-compose.yml:**
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=password
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: mydb
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"

volumes:
  mysql_data:
```

**运行:**
```bash
docker-compose up -d
```

---

## 总结

通过这些进阶功能，你可以将基础的 ViewSet 框架扩展为一个功能完整的企业级 API 框架，包括：

- ✅ 配置管理
- ✅ JWT 认证
- ✅ 权限控制 (RBAC)
- ✅ 数据验证
- ✅ 结构化日志
- ✅ API 文档
- ✅ Redis 缓存
- ✅ 数据库迁移
- ✅ 单元测试
- ✅ Docker 部署

这些功能都是模块化的，你可以根据项目需求选择性地集成。
