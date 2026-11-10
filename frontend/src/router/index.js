import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 路由表：页面在 DefaultLayout 内渲染（App.vue 承载布局）
const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/article/:id',
    name: 'article',
    component: () => import('@/views/ArticleDetail.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/Register.vue')
  },
  {
    path: '/editor',
    name: 'editor-create',
    component: () => import('@/views/Editor.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/editor/:id',
    name: 'editor-edit',
    component: () => import('@/views/Editor.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/me',
    name: 'me',
    component: () => import('@/views/Me.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/views/Admin.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/archive',
    name: 'archive',
    component: () => import('@/views/Archive.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// 全局前置守卫：登录 / 管理员鉴权（前端仅做体验层拦截，权限以服务端为准）
router.beforeEach((to) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLogin) {
    ElMessage.warning('请先登录')
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresAdmin) {
    if (!userStore.isLogin) {
      ElMessage.warning('请先登录')
      return { path: '/login', query: { redirect: to.fullPath } }
    }
    if (!userStore.isAdmin) {
      ElMessage.error('无权访问后台管理')
      return { path: '/' }
    }
  }

  return true
})

export default router
