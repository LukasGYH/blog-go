import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useUserStore } from '@/stores/user'

// 统一构造可被调用方识别的业务错误对象，便于页面做降级处理
function createApiError(code, message, detail) {
  const err = new Error(message || '请求失败')
  err.code = code
  err.detail = detail
  return err
}

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

// 请求拦截器：有 token 则注入 Authorization
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('blog_token')
    if (token) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器：统一剥壳 { code, message, data }
request.interceptors.response.use(
  (response) => {
    const body = response.data
    // 非标准响应体，直接返回
    if (!body || typeof body !== 'object' || !('code' in body)) {
      return body
    }
    // 业务成功
    if (body.code === 200 || body.code === 201) {
      return body.data
    }
    // 后端以 200 返回但业务码非成功（含 501 待实现）
    if (body.code === 501) {
      ElMessage.warning(body.message || '该功能后端待实现')
      return Promise.reject(createApiError(501, body.message || '该功能后端待实现'))
    }
    ElMessage.error(body.message || '请求失败')
    return Promise.reject(createApiError(body.code, body.message, body.error))
  },
  (error) => {
    const { response } = error

    // 网络错误 / 后端未启动
    if (!response) {
      const msg = '后端服务未启动，请先运行 make run'
      ElMessage.error(msg)
      return Promise.reject(createApiError(-1, msg))
    }

    const status = response.status
    const body = response.data || {}
    const message = body.message || body.error || `请求失败（${status}）`

    // 未认证：清空登录态并跳转登录
    if (status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      ElMessage.error('登录已过期，请重新登录')
      const current = router.currentRoute.value
      if (current.path !== '/login') {
        router.push({ path: '/login', query: { redirect: current.fullPath } })
      }
      return Promise.reject(createApiError(status, message))
    }

    // 501：后端功能待实现，不跳登录、不崩溃，仅提示并 reject
    if (status === 501) {
      const msg = body.message || '该功能后端待实现'
      ElMessage.warning(msg)
      return Promise.reject(createApiError(501, msg))
    }

    ElMessage.error(message)
    return Promise.reject(createApiError(status, message, body.error))
  }
)

export default request
