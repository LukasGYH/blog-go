<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listArticles, removeArticle, changeArticleStatus } from '@/api/article'
import { updateMe, updatePassword } from '@/api/auth'
import { useUserStore } from '@/stores/user'
import { formatDateTime } from '@/utils/format'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('articles')

/* ============ 我的文章 ============ */
const myArticles = ref([])
const articleTotal = ref(0)
const articlePage = ref(1)
const articlePageSize = ref(10)
const statusFilter = ref('')
const articlesLoading = ref(false)
const articlesUnimplemented = ref(false)

async function fetchMyArticles() {
  articlesLoading.value = true
  try {
    const data = await listArticles({
      page: articlePage.value,
      page_size: articlePageSize.value,
      author_id: userStore.user?.id,
      status: statusFilter.value || undefined
    })
    myArticles.value = data?.list || []
    articleTotal.value = data?.total || 0
    articlesUnimplemented.value = false
  } catch (err) {
    myArticles.value = []
    articleTotal.value = 0
    if (err?.code === 501) articlesUnimplemented.value = true
  } finally {
    articlesLoading.value = false
  }
}

function handleArticleFilter() {
  articlePage.value = 1
  fetchMyArticles()
}

function handleArticlePage(p) {
  articlePage.value = p
  fetchMyArticles()
}

function goEdit(id) {
  router.push(`/editor/${id}`)
}

async function handleDeleteArticle(row) {
  try {
    await ElMessageBox.confirm(`确定删除《${row.title}》吗？`, '提示', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消'
    })
  } catch {
    return
  }
  try {
    await removeArticle(row.id)
    ElMessage.success('已删除')
    fetchMyArticles()
  } catch {
    // 错误提示已由响应拦截器统一处理
  }
}

// 切换草稿/发布状态
async function toggleStatus(row) {
  const next = row.status === 'published' ? 'draft' : 'published'
  try {
    await changeArticleStatus(row.id, next)
    row.status = next
    ElMessage.success(next === 'published' ? '已发布' : '已转为草稿')
  } catch {
    // 错误提示已由响应拦截器统一处理
  }
}

/* ============ 个人资料 ============ */
const profileForm = reactive({
  nickname: '',
  email: '',
  avatar: '',
  bio: ''
})
const profileSaving = ref(false)

function fillProfile() {
  const u = userStore.user || {}
  profileForm.nickname = u.nickname || ''
  profileForm.email = u.email || ''
  profileForm.avatar = u.avatar || ''
  profileForm.bio = u.bio || ''
}

async function saveProfile() {
  profileSaving.value = true
  try {
    const data = await updateMe({ ...profileForm })
    // 后端可能返回更新后的用户信息，缺失时本地合并
    userStore.setUser(data && data.id ? data : { ...userStore.user, ...profileForm })
    ElMessage.success('资料已更新')
  } catch {
    // 错误提示已由响应拦截器统一处理
  } finally {
    profileSaving.value = false
  }
}

/* ============ 修改密码 ============ */
const passwordFormRef = ref()
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})
const passwordSaving = ref(false)

function validateConfirm(rule, value, callback) {
  if (value !== passwordForm.new_password) callback(new Error('两次输入的新密码不一致'))
  else callback()
}

const passwordRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 32, message: '密码长度为 6-32 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' }
  ]
}

async function savePassword() {
  if (!passwordFormRef.value) return
  try {
    await passwordFormRef.value.validate()
  } catch {
    return
  }
  passwordSaving.value = true
  try {
    await updatePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    })
    ElMessage.success('密码修改成功')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch {
    // 错误提示已由响应拦截器统一处理
  } finally {
    passwordSaving.value = false
  }
}

onMounted(() => {
  fillProfile()
  fetchMyArticles()
})
</script>

<template>
  <el-card shadow="never" class="me-card">
    <el-tabs v-model="activeTab">
      <!-- 我的文章 -->
      <el-tab-pane label="我的文章" name="articles">
        <div class="tab-toolbar">
          <el-radio-group v-model="statusFilter" @change="handleArticleFilter">
            <el-radio-button value="">全部</el-radio-button>
            <el-radio-button value="published">已发布</el-radio-button>
            <el-radio-button value="draft">草稿</el-radio-button>
          </el-radio-group>
          <el-button type="primary" @click="router.push('/editor')">
            <el-icon><EditPen /></el-icon>
            <span>写新文章</span>
          </el-button>
        </div>

        <el-table
          v-loading="articlesLoading"
          :data="myArticles"
          class="article-table"
          empty-text="暂无文章"
        >
          <el-table-column prop="title" label="标题" min-width="220">
            <template #default="{ row }">
              <router-link :to="`/article/${row.id}`" class="table-title">
                {{ row.title }}
              </router-link>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
                {{ row.status === 'published' ? '已发布' : '草稿' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="view_count" label="浏览" width="90" />
          <el-table-column label="发布时间" width="170">
            <template #default="{ row }">
              {{ formatDateTime(row.published_at || row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="240" fixed="right">
            <template #default="{ row }">
              <el-button text size="small" @click="goEdit(row.id)">编辑</el-button>
              <el-button text size="small" @click="toggleStatus(row)">
                {{ row.status === 'published' ? '转草稿' : '发布' }}
              </el-button>
              <el-button text size="small" type="danger" @click="handleDeleteArticle(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-empty
          v-if="!articlesLoading && !myArticles.length && articlesUnimplemented"
          description="文章列表暂不可用：后端 ArticleService 待实现"
        />

        <div v-if="articleTotal > articlePageSize" class="pagination">
          <el-pagination
            :current-page="articlePage"
            :page-size="articlePageSize"
            :total="articleTotal"
            layout="prev, pager, next"
            background
            @current-change="handleArticlePage"
          />
        </div>
      </el-tab-pane>

      <!-- 个人资料 -->
      <el-tab-pane label="个人资料" name="profile">
        <el-form :model="profileForm" label-width="90px" class="narrow-form">
          <el-form-item label="昵称">
            <el-input v-model="profileForm.nickname" placeholder="昵称" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="profileForm.email" placeholder="邮箱" />
          </el-form-item>
          <el-form-item label="头像 URL">
            <el-input v-model="profileForm.avatar" placeholder="头像图片 URL" />
          </el-form-item>
          <el-form-item label="个人简介">
            <el-input
              v-model="profileForm.bio"
              type="textarea"
              :rows="3"
              maxlength="255"
              show-word-limit
              placeholder="一句话介绍自己"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="profileSaving" @click="saveProfile">
              保存资料
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 修改密码 -->
      <el-tab-pane label="修改密码" name="password">
        <el-form
          ref="passwordFormRef"
          :model="passwordForm"
          :rules="passwordRules"
          label-width="90px"
          class="narrow-form"
        >
          <el-form-item label="当前密码" prop="old_password">
            <el-input
              v-model="passwordForm.old_password"
              type="password"
              show-password
              placeholder="请输入当前密码"
            />
          </el-form-item>
          <el-form-item label="新密码" prop="new_password">
            <el-input
              v-model="passwordForm.new_password"
              type="password"
              show-password
              placeholder="6-32 位"
            />
          </el-form-item>
          <el-form-item label="确认新密码" prop="confirm_password">
            <el-input
              v-model="passwordForm.confirm_password"
              type="password"
              show-password
              placeholder="请再次输入新密码"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="passwordSaving" @click="savePassword">
              修改密码
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<style scoped>
.me-card {
  border-radius: var(--radius);
}
.tab-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.article-table {
  width: 100%;
}
.table-title {
  color: var(--text-main);
}
.table-title:hover {
  color: var(--primary);
}
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}
.narrow-form {
  max-width: 520px;
  padding-top: 8px;
}
</style>
