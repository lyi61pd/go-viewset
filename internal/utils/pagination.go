package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
	Limit    int
}

// GetPaginationParams 从 gin.Context 中获取分页参数
// 支持两种方式：
// 1. page + page_size
// 2. limit + offset
func GetPaginationParams(c *gin.Context) *PaginationParams {
	params := &PaginationParams{
		Page:     1,
		PageSize: 10,
	}

	// 优先使用 page + page_size
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			params.PageSize = pageSize
			// 限制最大 page_size
			if params.PageSize > 100 {
				params.PageSize = 100
			}
		}
	}

	// 如果提供了 limit 和 offset，则使用这些参数
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
			if params.Limit > 100 {
				params.Limit = 100
			}
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}

	// 如果使用了 limit/offset，则计算对应的 page/page_size
	if params.Limit > 0 {
		params.PageSize = params.Limit
		if params.Offset > 0 {
			params.Page = (params.Offset / params.PageSize) + 1
		}
	} else {
		// 否则根据 page/page_size 计算 offset/limit
		params.Offset = (params.Page - 1) * params.PageSize
		params.Limit = params.PageSize
	}

	return params
}

// ApplyPagination 对 GORM 查询应用分页
func ApplyPagination(db *gorm.DB, params *PaginationParams) *gorm.DB {
	return db.Offset(params.Offset).Limit(params.Limit)
}

// GetTotal 获取总记录数
func GetTotal(db *gorm.DB, model interface{}) int64 {
	var total int64
	db.Model(model).Count(&total)
	return total
}

// BuildPagination 构建分页信息
func BuildPagination(params *PaginationParams, total int64) *Pagination {
	return &Pagination{
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}
}
