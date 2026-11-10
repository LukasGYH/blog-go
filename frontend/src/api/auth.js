import request from './request'

// 认证相关接口

// 注册
export function register(data) {
  return request.post('/auth/register', data)
}

// 登录
export function login(data) {
  return request.post('/auth/login', data)
}

// 当前登录用户
export function me() {
  return request.get('/auth/me')
}

// 修改个人资料
export function updateMe(data) {
  return request.put('/auth/me', data)
}

// 修改密码
export function updatePassword(data) {
  return request.put('/auth/password', data)
}
