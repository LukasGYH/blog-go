<script setup>
import { ref, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listComments, createComment, removeComment } from '@/api/comment'
import { useUserStore } from '@/stores/user'
import CommentItem from './CommentItem.vue'

const props = defineProps({
  articleId: { type: [Number, String], required: true },
  articleAuthorId: { type: [Number, String], default: null }
})

const router = useRouter()
const userStore = useUserStore()

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
// 后端 CommentService 尚未实现（501）时标记，用于优雅降级
const unimplemented = ref(false)

const content = ref('')
const replyTo = ref(null)
const submitting = ref(false)

const placeholder = computed(() =>
  replyTo.value ? `回复 @${nameOf(replyTo.value.user)}：` : '写下你的评论…'
)

function nameOf(user) {
  return user?.nickname || user?.username || '匿名'
}

async function fetchList() {
  loading.value = true
  try {
    const data = await listComments(props.articleId, {
      page: page.value,
      page_size: pageSize.value
    })
    list.value = data?.list || []
    total.value = data?.total || 0
    unimplemented.value = false
  } catch (err) {
    // 501 等后端待实现场景：空状态兜底，不崩溃
    list.value = []
    total.value = 0
    if (err?.code === 501) unimplemented.value = true
  } finally {
    loading.value = false
  }
}

function handleReply(comment) {
  replyTo.value = comment
}

function cancelReply() {
  replyTo.value = null
}

async function submit() {
  const text = content.value.trim()
  if (!text) {
    ElMessage.warning('评论内容不能为空')
    return
  }
  if (text.length > 500) {
    ElMessage.warning('评论内容不能超过 500 字')
    return
  }

  let parentId = 0
  if (replyTo.value) {
    // 若回复对象本身是回复，则归并到其顶层父评论
    parentId = replyTo.value.parent_id ? replyTo.value.parent_id : replyTo.value.id
  }

  submitting.value = true
  try {
    await createComment(props.articleId, { content: text, parent_id: parentId })
    ElMessage.success('评论成功')
    content.value = ''
    replyTo.value = null
    page.value = 1
    await fetchList()
  } catch {
    // 错误提示已由响应拦截器统一处理
  } finally {
    submitting.value = false
  }
}

async function handleRemove(comment) {
  try {
    await ElMessageBox.confirm('确定删除这条评论吗？', '提示', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消'
    })
  } catch {
    return
  }
  try {
    await removeComment(comment.id)
    ElMessage.success('已删除')
    await fetchList()
  } catch {
    // 错误提示已由响应拦截器统一处理
  }
}

function goLogin() {
  router.push({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } })
}

watch(
  () => props.articleId,
  () => {
    page.value = 1
    fetchList()
  },
  { immediate: true }
)

watch(page, () => fetchList())
watch(pageSize, () => {
  page.value = 1
})
</script>

<template>
  <el-card class="comment-section" shadow="never">
    <template #header>
      <div class="comment-header">
        <span class="card-title" style="margin: 0">评论</span>
        <span class="comment-count">共 {{ total }} 条</span>
      </div>
    </template>

    <!-- 发表评论 -->
    <div v-if="userStore.isLogin" class="comment-editor">
      <el-input
        v-model="content"
        type="textarea"
        :rows="3"
        :placeholder="placeholder"
        maxlength="500"
        show-word-limit
      />
      <div class="editor-actions">
        <el-button v-if="replyTo" text @click="cancelReply">取消回复</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">发表</el-button>
      </div>
    </div>
    <div v-else class="comment-login-tip">
      <el-button type="primary" plain @click="goLogin">登录后参与评论</el-button>
    </div>

    <!-- 评论列表 -->
    <div v-loading="loading" class="comment-list">
      <template v-if="list.length">
        <CommentItem
          v-for="item in list"
          :key="item.id"
          :comment="item"
          :article-author-id="articleAuthorId"
          @reply="handleReply"
          @remove="handleRemove"
        />
      </template>
      <el-empty
        v-else-if="!loading"
        :description="unimplemented ? '评论功能后端 CommentService 待实现' : '还没有评论，来抢沙发吧'"
      />
    </div>

    <div v-if="total > pageSize" class="comment-pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="prev, pager, next"
        background
        small
      />
    </div>
  </el-card>
</template>

<style scoped>
.comment-section {
  border-radius: var(--radius);
  margin-top: 16px;
}
.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.comment-count {
  color: var(--text-muted);
  font-size: 13px;
}
.comment-editor {
  margin-bottom: 8px;
}
.editor-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 10px;
}
.comment-login-tip {
  padding: 16px;
  background: #f7f8fa;
  border-radius: var(--radius);
  text-align: center;
  margin-bottom: 8px;
}
.comment-list {
  min-height: 60px;
}
.comment-pagination {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}
</style>
