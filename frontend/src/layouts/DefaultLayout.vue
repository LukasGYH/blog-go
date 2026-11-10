<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { initialOf } from '@/utils/format'

const router = useRouter()
const userStore = useUserStore()

const navItems = [
  { label: '首页', path: '/' },
  { label: '归档', path: '/archive' }
]

const avatarLetter = computed(() => initialOf(userStore.displayName))

// 跳转到新建文章
function goEditor() {
  router.push('/editor')
}

// 下拉菜单操作
function handleCommand(command) {
  if (command === 'me') router.push('/me')
  else if (command === 'admin') router.push('/admin')
  else if (command === 'editor') router.push('/editor')
  else if (command === 'logout') {
    userStore.logout()
    ElMessage.success('已退出登录')
    router.push('/')
  }
}
</script>

<template>
  <div class="layout">
    <header class="layout-header">
      <div class="app-container header-inner">
        <router-link to="/" class="brand">
          <span class="brand-dot"></span>
          <span class="brand-name">Lukas 的博客</span>
        </router-link>

        <nav class="nav">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="nav-link"
            active-class="nav-link--active"
          >
            {{ item.label }}
          </router-link>
        </nav>

        <div class="header-actions">
          <template v-if="userStore.isLogin">
            <el-button type="primary" @click="goEditor">
              <el-icon><EditPen /></el-icon>
              <span>写文章</span>
            </el-button>
            <el-dropdown @command="handleCommand">
              <span class="user-chip">
                <span v-if="userStore.user?.avatar" class="user-avatar">
                  <img :src="userStore.user.avatar" alt="avatar" />
                </span>
                <span v-else class="avatar-letter user-avatar-letter">
                  {{ avatarLetter }}
                </span>
                <span class="user-name">{{ userStore.displayName }}</span>
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="me">个人中心</el-dropdown-item>
                  <el-dropdown-item v-if="userStore.isAdmin" command="admin">
                    后台管理
                  </el-dropdown-item>
                  <el-dropdown-item command="editor">写文章</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button text @click="router.push('/login')">登录</el-button>
            <el-button type="primary" @click="router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </header>

    <main class="layout-main">
      <div class="app-container">
        <router-view />
      </div>
    </main>

    <footer class="layout-footer">
      <div class="app-container footer-inner">
        <span>© {{ new Date().getFullYear() }} Lukas 的博客 · Powered by Vue 3 + Vite + Gin</span>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.layout-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  border-bottom: 1px solid var(--border-light);
}

.header-inner {
  height: 60px;
  display: flex;
  align-items: center;
  gap: 24px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
  color: var(--text-main);
}
.brand-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--primary);
}

.nav {
  display: flex;
  gap: 8px;
  flex: 1;
}
.nav-link {
  padding: 6px 12px;
  border-radius: 6px;
  color: var(--text-sub);
}
.nav-link:hover {
  background: #f2f6fc;
  color: var(--primary);
}
.nav-link--active {
  color: var(--primary);
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  outline: none;
  color: var(--text-main);
}
.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
}
.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.user-avatar-letter {
  width: 32px;
  height: 32px;
  font-size: 14px;
}
.user-name {
  font-size: 14px;
}

.layout-main {
  flex: 1;
  padding: 24px 0 40px;
}

.layout-footer {
  border-top: 1px solid var(--border-light);
  background: #fff;
}
.footer-inner {
  padding: 20px 16px;
  text-align: center;
  color: var(--text-muted);
  font-size: 13px;
}
</style>
