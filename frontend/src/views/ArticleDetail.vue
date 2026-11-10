<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getArticle, removeArticle } from '@/api/article'
import { useUserStore } from '@/stores/user'
import { formatDateTime, initialOf } from '@/utils/format'
import MarkdownView from '@/components/MarkdownView.vue'
import CommentList from '@/components/CommentList.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref(null)
const loading = ref(false)
const notFound = ref(false)
const unimplemented = ref(false)

const authorName = computed(
  () => article.value?.author?.nickname || article.value?.author?.username || '匿名'
)
const avatarLetter = computed(() => initialOf(authorName.value))
const tags = computed(() => article.value?.tags || [])

// 作者或管理员可编辑/删除（前端仅体验层显示，权限以服务端为准）
const canManage = computed(() => {
  if (!userStore.isLogin || !article.value) return false
  if (userStore.isAdmin) return true
  const uid = userStore.user?.id
  return article.value.author?.id != null && article.value.author.id === uid
})

async function fetchDetail() {
  loading.value = true
  notFound.value = false
  unimplemented.value = false
  try {
    article.value = await getArticle(route.params.id)
  } catch (err) {
    article.value = null
    if (err?.code === 501) unimplemented.value = true
    else if (err?.code === 404) notFound.value = true
  } finally {
    loading.value = false
  }
}

function goEditor() {
  router.push(`/editor/${article.value.id}`)
}

async function handleDelete() {
  try {
    await ElMessageBox.confirm('确定删除这篇文章吗？删除后可在数据库保留软删除记录。', '提示', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消'
    })
  } catch {
    return
  }
  try {
    await removeArticle(article.value.id)
    ElMessage.success('已删除')
    router.push('/')
  } catch {
    // 错误提示已由响应拦截器统一处理
  }
}

watch(
  () => route.params.id,
  () => fetchDetail(),
  { immediate: true }
)
</script>

<template>
  <div class="article-detail">
    <el-card v-loading="loading" shadow="never" class="detail-card">
      <template v-if="article">
        <h1 class="detail-title">{{ article.title }}</h1>

        <div class="detail-meta">
          <span class="meta-author">
            <span v-if="article.author?.avatar" class="author-avatar">
              <img :src="article.author.avatar" alt="avatar" />
            </span>
            <span v-else class="avatar-letter author-avatar-letter">
              {{ avatarLetter }}
            </span>
            {{ authorName }}
          </span>
          <span class="meta-item">
            <el-icon><Clock /></el-icon>
            {{ formatDateTime(article.published_at || article.created_at) }}
          </span>
          <span class="meta-item">
            <el-icon><View /></el-icon>{{ article.view_count ?? 0 }} 次浏览
          </span>
          <el-tag v-if="article.category" size="small" type="info" effect="plain">
            {{ article.category.name }}
          </el-tag>
          <el-tag
            v-for="tag in tags"
            :key="tag.id"
            size="small"
            effect="plain"
          >
            {{ tag.name }}
          </el-tag>
          <el-tag v-if="article.status === 'draft'" size="small" type="warning">
            草稿
          </el-tag>

          <div v-if="canManage" class="meta-actions">
            <el-button size="small" @click="goEditor">编辑</el-button>
            <el-button size="small" type="danger" plain @click="handleDelete">
              删除
            </el-button>
          </div>
        </div>

        <el-divider />

        <MarkdownView :content="article.content" />
      </template>

      <el-empty
        v-else-if="!loading && notFound"
        description="文章不存在或已被删除"
      />
      <el-empty
        v-else-if="!loading && unimplemented"
        description="文章详情暂不可用：后端 ArticleService 待实现"
      />
      <el-empty v-else-if="!loading" description="加载失败" />
    </el-card>

    <CommentList
      v-if="article"
      :article-id="article.id"
      :article-author-id="article.author?.id"
    />
  </div>
</template>

<style scoped>
.detail-card {
  border-radius: var(--radius);
}
.detail-title {
  margin: 0 0 16px;
  font-size: 26px;
  font-weight: 700;
  line-height: 1.4;
}
.detail-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 14px;
  color: var(--text-muted);
  font-size: 13px;
}
.meta-author {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-sub);
}
.author-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  overflow: hidden;
}
.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.author-avatar-letter {
  width: 28px;
  height: 28px;
  font-size: 13px;
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.meta-actions {
  margin-left: auto;
  display: flex;
  gap: 8px;
}
</style>
