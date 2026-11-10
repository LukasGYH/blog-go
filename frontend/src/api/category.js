import request from './request'

// 分类相关接口

// 分类列表（含 article_count）
export function listCategories() {
  return request.get('/categories')
}

// 新增分类（管理员）
export function createCategory(data) {
  return request.post('/categories', data)
}

// 修改分类（管理员）
export function updateCategory(id, data) {
  return request.put(`/categories/${id}`, data)
}

// 删除分类（管理员，有文章返回 409）
export function removeCategory(id) {
  return request.delete(`/categories/${id}`)
}
