package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FilterParams 过滤参数
type FilterParams struct {
	Filters  map[string]interface{}
	OrderBy  string
	OrderDir string
}

// GetFilterParams 从 gin.Context 中获取过滤参数
// 支持：
// 1. 简单的等值过滤：?name=abc&status=active
// 2. 排序：?order_by=created_at desc 或 ?ordering=-created_at
func GetFilterParams(c *gin.Context, excludeKeys ...string) *FilterParams {
	params := &FilterParams{
		Filters: make(map[string]interface{}),
	}

	// 需要排除的特殊参数
	excludeMap := map[string]bool{
		"page":      true,
		"page_size": true,
		"limit":     true,
		"offset":    true,
		"order_by":  true,
		"ordering":  true,
	}

	// 添加用户自定义的排除参数
	for _, key := range excludeKeys {
		excludeMap[key] = true
	}

	// 获取所有查询参数
	for key, values := range c.Request.URL.Query() {
		if excludeMap[key] {
			continue
		}
		if len(values) > 0 {
			// 如果有多个值，取第一个
			params.Filters[key] = values[0]
		}
	}

	// 处理排序参数
	// 支持两种格式：
	// 1. order_by=created_at desc
	// 2. ordering=-created_at (DRF 风格)
	if orderBy := c.Query("order_by"); orderBy != "" {
		parts := strings.Fields(orderBy)
		if len(parts) > 0 {
			params.OrderBy = parts[0]
			if len(parts) > 1 {
				params.OrderDir = strings.ToUpper(parts[1])
			} else {
				params.OrderDir = "ASC"
			}
		}
	} else if ordering := c.Query("ordering"); ordering != "" {
		if strings.HasPrefix(ordering, "-") {
			params.OrderBy = ordering[1:]
			params.OrderDir = "DESC"
		} else {
			params.OrderBy = ordering
			params.OrderDir = "ASC"
		}
	}

	// 验证排序方向
	if params.OrderDir != "" && params.OrderDir != "ASC" && params.OrderDir != "DESC" {
		params.OrderDir = "ASC"
	}

	return params
}

// ApplyFilters 对 GORM 查询应用过滤
func ApplyFilters(db *gorm.DB, params *FilterParams) *gorm.DB {
	// 应用等值过滤
	for key, value := range params.Filters {
		// 使用参数化查询防止 SQL 注入
		db = db.Where(key+" = ?", value)
	}

	// 应用排序
	if params.OrderBy != "" {
		// 验证字段名，防止 SQL 注入
		// 这里简单处理，实际项目中应该有白名单验证
		orderClause := sanitizeOrderBy(params.OrderBy)
		if params.OrderDir != "" {
			orderClause += " " + params.OrderDir
		}
		db = db.Order(orderClause)
	}

	return db
}

// sanitizeOrderBy 清理排序字段名，防止 SQL 注入
func sanitizeOrderBy(field string) string {
	// 移除危险字符，只保留字母、数字、下划线和点
	var result strings.Builder
	for _, r := range field {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '.' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ApplySearch 应用模糊搜索（可选功能）
// 使用方式：?search=keyword
func ApplySearch(db *gorm.DB, c *gin.Context, fields ...string) *gorm.DB {
	search := c.Query("search")
	if search == "" || len(fields) == 0 {
		return db
	}

	// 构建 OR 查询
	query := db
	for i, field := range fields {
		if i == 0 {
			query = query.Where(field+" LIKE ?", "%"+search+"%")
		} else {
			query = query.Or(field+" LIKE ?", "%"+search+"%")
		}
	}

	return query
}
