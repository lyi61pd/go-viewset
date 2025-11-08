package main

import (
	"fmt"
	"go-viewset/internal/config"
	"go-viewset/internal/models"
	"go-viewset/internal/router"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置路由
	r := router.SetupRouter(db)

	// 启动服务
	fmt.Printf("🚀 服务启动成功，监听端口: %s\n", cfg.Server.Port)
	fmt.Println("📚 API 文档:")
	fmt.Println("  - GET    /api/users/          获取用户列表")
	fmt.Println("  - GET    /api/users/:id       获取单个用户")
	fmt.Println("  - POST   /api/users/          创建用户")
	fmt.Println("  - PUT    /api/users/:id       更新用户")
	fmt.Println("  - DELETE /api/users/:id       删除用户")
	fmt.Println("  - POST   /api/users/:id/activate      激活用户")
	fmt.Println("  - POST   /api/users/:id/deactivate    停用用户")
	fmt.Println("  - POST   /api/users/:id/reset_password 重置密码")
	fmt.Println("  - GET    /api/users/stats     获取统计信息")
	fmt.Println("")

	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// initDB 初始化数据库
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 构建 DSN 连接字符串
	dsn := cfg.Database.GetDSN()

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 创建一些示例数据
	createSampleData(db)

	return db, nil
}

// createSampleData 创建示例数据
func createSampleData(db *gorm.DB) {
	// 检查是否已有数据
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	// 创建示例用户
	users := []models.User{
		{
			Name:   "张三",
			Email:  "zhangsan@example.com",
			Status: "active",
			Age:    25,
			Phone:  "13800138000",
		},
		{
			Name:   "李四",
			Email:  "lisi@example.com",
			Status: "active",
			Age:    30,
			Phone:  "13800138001",
		},
		{
			Name:   "王五",
			Email:  "wangwu@example.com",
			Status: "inactive",
			Age:    28,
			Phone:  "13800138002",
		},
	}

	for _, user := range users {
		db.Create(&user)
	}

	fmt.Println("✅ 示例数据创建成功")
}
