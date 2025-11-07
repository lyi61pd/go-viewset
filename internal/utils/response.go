package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination 分页信息
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithPagination 带分页的成功响应
func SuccessWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, Response{
		Code:       0,
		Msg:        "success",
		Data:       data,
		Pagination: pagination,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// ErrorWithStatus 带 HTTP 状态码的错误响应
func ErrorWithStatus(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusBadRequest, http.StatusBadRequest, msg)
}

// NotFound 404 错误
func NotFound(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusNotFound, http.StatusNotFound, msg)
}

// InternalServerError 500 错误
func InternalServerError(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, msg)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusUnauthorized, http.StatusUnauthorized, msg)
}

// Forbidden 403 错误
func Forbidden(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusForbidden, http.StatusForbidden, msg)
}
