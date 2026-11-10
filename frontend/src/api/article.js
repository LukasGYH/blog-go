import request from './request'

// 文章相关接口

// 文章列表（分页 + 筛选）
export function listArticles(params) {
  return request.get('/articles', { params })
}

// 文章详情
export function getArticle(id) {
  return request.get(`/articles/${id}`)
}

// 新建文章
export function createArticle(data) {
  return request.post('/articles', data)
}

// 编辑文章
export function updateArticle(id, data) {
  return request.put(`/articles/${id}`, data)
}

// 删除文章（软删除）
export function removeArticle(id) {
  return request.delete(`/articles/${id}`)
}

// 切换文章状态（草稿/发布）
export function changeArticleStatus(id, status) {
  return request.patch(`/articles/${id}/status`, { status })
}

// 归档统计
export function getArchive() {
  return request.get('/articles/archive')
}

// 热门文章
export function getHotArticles(limit = 5) {
  return request.get('/articles/hot', { params: { limit } })
}
