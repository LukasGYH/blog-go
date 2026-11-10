import request from './request'

// 评论相关接口

// 评论列表（该文章的顶层评论，内嵌 replies）
export function listComments(articleId, params) {
  return request.get(`/articles/${articleId}/comments`, { params })
}

// 发表评论 / 回复（parent_id 为 0 或省略表示顶层评论）
export function createComment(articleId, data) {
  return request.post(`/articles/${articleId}/comments`, data)
}

// 删除评论
export function removeComment(id) {
  return request.delete(`/comments/${id}`)
}
