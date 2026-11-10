package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError 业务错误：携带 HTTP 状态码与对外提示文案
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 支持 errors.Is / errors.As 追溯底层错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建业务错误
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// BadRequest 400 参数错误
func BadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message)
}

// Unauthorized 401 未认证
func Unauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message)
}

// Forbidden 403 无权限
func Forbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message)
}

// NotFound 404 资源不存在
func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message)
}

// Conflict 409 资源冲突
func Conflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message)
}

// NotImplemented 501 功能待实现，供骨架中的空实现统一返回
func NotImplemented(method string) *AppError {
	return NewAppError(http.StatusNotImplemented, "功能待实现："+method)
}

// HandleError 统一错误出口：识别 AppError 直接返回；
// 未识别错误记日志并返回 500，不向客户端泄露内部信息
func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.Code, appErr.Message)
		return
	}
	log.Printf("[error] %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	Error(c, http.StatusInternalServerError, "internal server error")
}
