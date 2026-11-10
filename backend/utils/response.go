package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response 统一响应体：{code, message, data, error}
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// PageData 统一分页数据体：{list, total, page, page_size}
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 返回 200 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Created 返回 201 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "success",
		Data:    data,
	})
}

// SuccessPage 返回 200 分页响应，list 需为已初始化的切片以避免 JSON 出现 null
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Error 返回错误响应，error 字段与 message 保持一致
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   message,
	})
}

// ValidationError 返回 422 参数校验错误，error 字段为「字段 -> 提示」映射
func ValidationError(c *gin.Context, fields map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Code:    422,
		Message: "validation failed",
		Error:   fields,
	})
}

// ParseBindingError 将 binding 校验错误转换为「字段 -> 中文提示」映射
func ParseBindingError(err error) map[string]string {
	result := make(map[string]string)
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, fieldErr := range validationErrs {
			result[strings.ToLower(fieldErr.Field())] = bindingMessage(fieldErr)
		}
		return result
	}
	result["general"] = err.Error()
	return result
}

// bindingMessage 将校验标签翻译为中文提示
func bindingMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "该字段必填"
	case "min":
		return "取值/长度小于最小限制 " + fieldErr.Param()
	case "max":
		return "取值/长度超过最大限制 " + fieldErr.Param()
	case "email":
		return "邮箱格式不正确"
	case "oneof":
		return "取值必须为: " + fieldErr.Param()
	default:
		return "校验失败: " + fieldErr.Tag()
	}
}
