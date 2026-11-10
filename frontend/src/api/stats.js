import request from './request'

// 统计相关接口（管理员）

// 概览统计
export function getOverview() {
  return request.get('/stats/overview')
}

// 分类维度统计
export function getCategoryStats() {
  return request.get('/stats/categories')
}
