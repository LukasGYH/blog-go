package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog-api/utils"
)

// parseIDParam 解析路径参数 :id，失败时写入 400 响应并返回 false
func parseIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(id), true
}
