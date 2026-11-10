import request from './request'

// 标签相关接口

// 标签列表（含 article_count，可按关键字搜索）
export function listTags(params) {
  return request.get('/tags', { params })
}

// 新增标签（管理员）
export function createTag(data) {
  return request.post('/tags', data)
}

// 修改标签（管理员）
export function updateTag(id, data) {
  return request.put(`/tags/${id}`, data)
}

// 删除标签（管理员）
export function removeTag(id) {
  return request.delete(`/tags/${id}`)
}
