package viewset
package viewset

import (
	"fmt"
	"go-viewset/internal/utils"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BaseViewSet 定义 ViewSet 的基础接口
type BaseViewSet interface {
	List(c *gin.Context)
	Retrieve(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	RegisterRoutes(group *gin.RouterGroup)
}

// GenericViewSet 通用 ViewSet 实现
// 提供标准的 CRUD 操作，支持分页、过滤、排序
type GenericViewSet struct {
	DB        *gorm.DB
	Model     interface{}
	ModelType reflect.Type
}

// NewGenericViewSet 创建一个新的 GenericViewSet
// model 参数应该是一个模型的指针，例如 &User{}
func NewGenericViewSet(db *gorm.DB, model interface{}) *GenericViewSet {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return &GenericViewSet{
		DB:        db,
		Model:     model,
		ModelType: modelType,
	}
}

// List 获取列表
// 支持分页、过滤和排序
// GET /items/?page=1&page_size=10&name=abc&order_by=created_at desc
func (v *GenericViewSet) List(c *gin.Context) {
	// 创建模型切片
	sliceType := reflect.SliceOf(reflect.PtrTo(v.ModelType))
	results := reflect.New(sliceType).Interface()

	// 获取分页参数
	paginationParams := utils.GetPaginationParams(c)

	// 获取过滤参数
	filterParams := utils.GetFilterParams(c)

	// 构建查询
	query := v.DB.Model(v.Model)

	// 应用过滤
	query = utils.ApplyFilters(query, filterParams)

	// 获取总数（在应用分页之前）
	var total int64
	query.Count(&total)

	// 应用分页
	query = utils.ApplyPagination(query, paginationParams)

	// 执行查询
	if err := query.Find(results).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("查询失败: %v", err))
		return
	}

	// 构建分页信息
	pagination := utils.BuildPagination(paginationParams, total)

	// 返回结果
	utils.SuccessWithPagination(c, results, pagination)
}

// Retrieve 获取单个对象
// GET /items/:id
func (v *GenericViewSet) Retrieve(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "缺少 ID 参数")
		return
	}

	// 创建模型实例
	result := reflect.New(v.ModelType).Interface()

	// 查询
	if err := v.DB.First(result, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "记录不存在")
		} else {
			utils.InternalServerError(c, fmt.Sprintf("查询失败: %v", err))
		}
		return
	}

	utils.Success(c, result)
}

// Create 创建新对象
// POST /items/
func (v *GenericViewSet) Create(c *gin.Context) {
	// 创建模型实例
	obj := reflect.New(v.ModelType).Interface()

	// 绑定请求数据
	if err := c.ShouldBindJSON(obj); err != nil {
		utils.BadRequest(c, fmt.Sprintf("请求数据格式错误: %v", err))
		return
	}

	// 创建记录
	if err := v.DB.Create(obj).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("创建失败: %v", err))
		return
	}

	utils.Success(c, obj)
}

// Update 更新对象
// PUT /items/:id
func (v *GenericViewSet) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "缺少 ID 参数")
		return
	}

	// 先查询是否存在
	existing := reflect.New(v.ModelType).Interface()
	if err := v.DB.First(existing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "记录不存在")
		} else {
			utils.InternalServerError(c, fmt.Sprintf("查询失败: %v", err))
		}
		return
	}

	// 绑定更新数据
	updates := reflect.New(v.ModelType).Interface()
	if err := c.ShouldBindJSON(updates); err != nil {
		utils.BadRequest(c, fmt.Sprintf("请求数据格式错误: %v", err))
		return
	}

	// 更新记录
	if err := v.DB.Model(existing).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("更新失败: %v", err))
		return
	}

	// 重新查询获取最新数据
	result := reflect.New(v.ModelType).Interface()
	v.DB.First(result, id)

	utils.Success(c, result)
}

// Delete 删除对象
// DELETE /items/:id
func (v *GenericViewSet) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, "缺少 ID 参数")
		return
	}

	// 创建模型实例
	obj := reflect.New(v.ModelType).Interface()

	// 先查询是否存在
	if err := v.DB.First(obj, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "记录不存在")
		} else {
			utils.InternalServerError(c, fmt.Sprintf("查询失败: %v", err))
		}
		return
	}

	// 删除记录
	if err := v.DB.Delete(obj).Error; err != nil {
		utils.InternalServerError(c, fmt.Sprintf("删除失败: %v", err))
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// RegisterRoutes 注册标准 RESTful 路由
// 子类可以覆盖此方法来添加自定义路由
func (v *GenericViewSet) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/", v.List)
	group.GET("/:id", v.Retrieve)
	group.POST("/", v.Create)
	group.PUT("/:id", v.Update)
	group.DELETE("/:id", v.Delete)
}

// RegisterAction 注册自定义 action
// method: HTTP 方法，例如 "POST", "GET"
// path: 路径，例如 "/:id/activate"
// handler: 处理函数
func (v *GenericViewSet) RegisterAction(group *gin.RouterGroup, method, path string, handler gin.HandlerFunc) {
	switch method {
	case "GET":
		group.GET(path, handler)
	case "POST":
		group.POST(path, handler)
	case "PUT":
		group.PUT(path, handler)
	case "DELETE":
		group.DELETE(path, handler)
	case "PATCH":
		group.PATCH(path, handler)
	default:
		group.Any(path, handler)
	}
}

// GetObjectOr404 获取对象，如果不存在则返回 404
// 这是一个辅助方法，用于在自定义 action 中快速获取对象
func (v *GenericViewSet) GetObjectOr404(c *gin.Context, id string) (interface{}, bool) {
	// 转换 ID
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.BadRequest(c, "无效的 ID")
		return nil, false
	}

	// 创建模型实例
	obj := reflect.New(v.ModelType).Interface()

	// 查询
	if err := v.DB.First(obj, idInt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "记录不存在")
		} else {
			utils.InternalServerError(c, fmt.Sprintf("查询失败: %v", err))
		}
		return nil, false
	}

	return obj, true
}

// PerformCreate 创建前的钩子，子类可以覆盖
func (v *GenericViewSet) PerformCreate(c *gin.Context, obj interface{}) error {
	return nil
}

// PerformUpdate 更新前的钩子，子类可以覆盖
func (v *GenericViewSet) PerformUpdate(c *gin.Context, obj interface{}) error {
	return nil
}

// PerformDestroy 删除前的钩子，子类可以覆盖
func (v *GenericViewSet) PerformDestroy(c *gin.Context, obj interface{}) error {
	return nil
}
