package services

import (
	"context"
	"fmt"

	"blog-api/models/vo"
	"blog-api/repositories"
)

// StatsService 统计业务逻辑
type StatsService struct {
	repo *repositories.StatsRepository
}

// NewStatsService 创建统计服务
func NewStatsService(repo *repositories.StatsRepository) *StatsService {
	return &StatsService{repo: repo}
}

// Overview 概览统计
func (s *StatsService) Overview(ctx context.Context) (*vo.StatsOverviewVO, error) {
	overview, err := s.repo.Overview(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询概览统计失败: %w", err)
	}
	return overview, nil
}

// CategoryStats 分类维度统计
func (s *StatsService) CategoryStats(ctx context.Context) ([]vo.CategoryStatsVO, error) {
	stats, err := s.repo.CountByCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询分类统计失败: %w", err)
	}
	return stats, nil
}
