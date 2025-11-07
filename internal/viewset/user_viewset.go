package viewset

import (
	"fmt"
	"go-viewset/internal/models"
	"go-viewset/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserViewSet 用户 ViewSet
// 通过嵌入 GenericViewSet 快速实现 CRUD
type UserViewSet struct {
	*GenericViewSet
}

// NewUserViewSet 创建用户 ViewSet
func NewUserViewSet(db *gorm.DB) *UserViewSet {
	return &UserViewSet{
		GenericViewSet: NewGenericViewSet(db, &models.User{}),
	}
}

// RegisterRoutes 注册路由
// 除了标准的 CRUD 路由外，还注册自定义 action
func (v *UserViewSet) RegisterRoutes(group *gin.RouterGroup) {
	// 注册标准 RESTful 路由（使用子类的方法）
	group.GET("/", v.List) // 使用覆盖后的 List 方法
	group.POST("/", v.Create) // 使用覆盖后的 Create 方法

	// 注册自定义 action
	// POST /users/:id/activate - 激活用户
	v.RegisterAction(group, "POST", "/:id/activate", v.Activate)

	// POST /users/:id/deactivate - 停用用户
	v.RegisterAction(group, "POST", "/:id/deactivate", v.Deactivate)

	// POST /users/:id/reset_password - 重置密码
	v.RegisterAction(group, "POST", "/:id/reset_password", v.ResetPassword)

	// GET /users/stats - 获取统计信息（不需要 ID 的 action）
	v.RegisterAction(group, "GET", "/stats", v.GetStats)
}

// Activate 激活用户
// POST /users/:id/activate
func (v *UserViewSet) Activate(c *gin.Context) {
	id := c.Param("id")

	// 使用辅助方法获取对象
	obj, ok := v.GetObjectOr404(c, id)
	if !ok {
		return
	}

	user := obj.(*models.User)

	// 更新状态
	user.Status = "active"
	if err := v.DB.Save(user).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("激活失败: %v", err))
		return
	}

	utils.Success(c, gin.H{
		"message": "用户已激活",
		"user":    user,
	})
}

// Deactivate 停用用户
// POST /users/:id/deactivate
func (v *UserViewSet) Deactivate(c *gin.Context) {
	id := c.Param("id")

	obj, ok := v.GetObjectOr404(c, id)
	if !ok {
		return
	}

	user := obj.(*models.User)

	// 更新状态
	user.Status = "inactive"
	if err := v.DB.Save(user).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("停用失败: %v", err))
		return
	}

	utils.Success(c, gin.H{
		"message": "用户已停用",
		"user":    user,
	})
}

// ResetPassword 重置密码
// POST /users/:id/reset_password
func (v *UserViewSet) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	obj, ok := v.GetObjectOr404(c, id)
	if !ok {
		return
	}

	user := obj.(*models.User)

	// 这里只是示例，实际项目中应该有密码重置逻辑
	// 例如发送邮件、生成临时密码等

	utils.Success(c, gin.H{
		"message": "密码重置邮件已发送",
		"user_id": user.ID,
		"email":   user.Email,
	})
}

// GetStats 获取用户统计信息
// GET /users/stats
func (v *UserViewSet) GetStats(c *gin.Context) {
	var total int64
	var activeCount int64
	var inactiveCount int64

	// 统计总数
	v.DB.Model(&models.User{}).Count(&total)

	// 统计活跃用户
	v.DB.Model(&models.User{}).Where("status = ?", "active").Count(&activeCount)

	// 统计非活跃用户
	v.DB.Model(&models.User{}).Where("status = ?", "inactive").Count(&inactiveCount)

	utils.Success(c, gin.H{
		"total":    total,
		"active":   activeCount,
		"inactive": inactiveCount,
	})
}

// 可以覆盖父类的方法来自定义行为
// 例如：在创建用户前进行额外的验证

// List 覆盖列表方法，添加 keyword 搜索功能
// 支持通过 ?keyword=xxx 对 name、email、phone 进行模糊搜索
func (v *UserViewSet) List(c *gin.Context) {
	// 创建结果切片
	var users []models.User

	// 获取分页参数
	paginationParams := utils.GetPaginationParams(c)

	// 获取过滤参数
	filterParams := utils.GetFilterParams(c, "keyword") // 排除 keyword，因为我们要单独处理

	// 构建查询
	query := v.DB.Model(&models.User{})

	// 处理 keyword 搜索（多字段模糊匹配）
	keyword := c.Query("keyword")
	if keyword != "" {
		// 使用 OR 条件对多个字段进行模糊搜索
		query = query.Where(
			"name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	// 应用其他过滤条件（如 status、age 等）
	query = utils.ApplyFilters(query, filterParams)

	// 获取总数（在应用分页之前）
	var total int64
	query.Count(&total)

	// 应用分页
	query = utils.ApplyPagination(query, paginationParams)

	// 执行查询
	if err := query.Find(&users).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("查询失败: %v", err))
		return
	}

	// 构建分页信息
	pagination := utils.BuildPagination(paginationParams, total)

	// 返回结果
	utils.SuccessWithPagination(c, users, pagination)
}

// Create 覆盖创建方法，添加自定义逻辑
func (v *UserViewSet) Create(c *gin.Context) {
	var user models.User

	// 绑定请求数据
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(c, fmt.Sprintf("请求数据格式错误: %v", err))
		return
	}

	// 自定义验证：检查邮箱是否已存在
	var count int64
	v.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		utils.BadRequest(c, "该邮箱已被注册")
		return
	}

	// 设置默认状态
	if user.Status == "" {
		user.Status = "inactive"
	}

	// 创建用户
	if err := v.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("创建失败: %v", err))
		return
	}

	utils.Success(c, user)
}
