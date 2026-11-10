package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-api/utils"
)

// Recovery panic 恢复中间件，统一返回 JSON 500
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("[panic] %s %s: %v", c.Request.Method, c.Request.URL.Path, recovered)
		utils.Error(c, http.StatusInternalServerError, "internal server error")
		c.Abort()
	})
}
