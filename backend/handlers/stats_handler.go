package handlers

import (
	"github.com/gin-gonic/gin"

	"blog-api/services"
	"blog-api/utils"
)

// StatsHandler 统计接口（后台，仅管理员）
type StatsHandler struct {
	statsService *services.StatsService
}

// NewStatsHandler 创建统计处理器
func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// Overview 概览统计
func (h *StatsHandler) Overview(c *gin.Context) {
	overview, err := h.statsService.Overview(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, overview)
}

// Categories 分类维度统计
func (h *StatsHandler) Categories(c *gin.Context) {
	stats, err := h.statsService.CategoryStats(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	utils.Success(c, stats)
}
