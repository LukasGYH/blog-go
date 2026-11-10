<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listCategories,
  createCategory,
  updateCategory,
  removeCategory
} from '@/api/category'
import { listTags, createTag, updateTag, removeTag } from '@/api/tag'
import { getOverview, getCategoryStats } from '@/api/stats'

const activeTab = ref('overview')

/* ============ 概览 ============ */
const overview = ref(null)
const overviewLoading = ref(false)
const overviewUnimplemented = ref(false)

const statCards = computed(() => {
  const o = overview.value || {}
  return [
    { label: '文章总数', value: o.article_count ?? 0, color: '#409eff' },
    { label: '已发布', value: o.published_count ?? 0, color: '#67c23a' },
    { label: '草稿', value: o.draft_count ?? 0, color: '#e6a23c' },
    { label: '用户数', value: o.user_count ?? 0, color: '#909399' },
    { label: '评论数', value: o.comment_count ?? 0, color: '#f56c6c' },
    { label: '分类数', value: o.category_count ?? 0, color: '#9254de' },
    { label: '标签数', value: o.tag_count ?? 0, color: '#13c2c2' }
  ]
})

const categoryStats = ref([])
const categoryStatsLoading = ref(false)
const categoryStatsUnimplemented = ref(false)

const maxCategoryCount = computed(() => {
  const arr = categoryStats.value.map((c) => c.article_count || 0)
  return arr.length ? Math.max(...arr) : 0
})

async function fetchOverview() {
  overviewLoading.value = true
  try {
    overview.value = await getOverview()
    overviewUnimplemented.value = false
  } catch (err) {
    overview.value = null
    if (err?.code === 501) overviewUnimplemented.value = true
  } finally {
    overviewLoading.value = false
  }
}

async function fetchCategoryStats() {
  categoryStatsLoading.value = true
  try {
    categoryStats.value = (await getCategoryStats()) || []
    categoryStatsUnimplemented.value = false
  } catch (err) {
    categoryStats.value = []
    if (err?.code === 501) categoryStatsUnimplemented.value = true
  } finally {
    categoryStatsLoading.value = false
  }
}

function progressPercent(count) {
  if (!maxCategoryCount.value) return 0
  return Math.round(((count || 0) / maxCategoryCount.value) * 100)
}

/* ============ 分类管理 ============ */
const categories = ref([])
const categoriesLoading = ref(false)

async function fetchCategories() {
  categoriesLoading.value = true
  try {
    categories.value = (await listCategories()) || []
  } catch {
    categories.value = []
  } finally {
    categoriesLoading.value = false
  }
}

const categoryDialogVisible = ref(false)
const categoryDialogTitle = ref('新增分类')
const categoryFormRef = ref()
const categoryForm = reactive({
  id: null,
  name: '',
  slug: '',
  description: '',
  sort_order: 0
})
const categoryRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

function openCategoryDialog(row) {
  if (row) {
    categoryDialogTitle.value = '编辑分类'
    categoryForm.id = row.id
    categoryForm.name = row.name
    categoryForm.slug = row.slug
    categoryForm.description = row.description || ''
    categoryForm.sort_order = row.sort_order ?? 0
  } else {
    categoryDialogTitle.value = '新增分类'
    categoryForm.id = null
    categoryForm.name = ''
    categoryForm.slug = ''
    categoryForm.description = ''
    categoryForm.sort_order = 0
  }
  categoryDialogVisible.value = true
}

async function submitCategory() {
  if (!categoryFormRef.value) return
  try {
    await categoryFormRef.value.validate()
  } catch {
    return
  }
  const payload = {
    name: categoryForm.name,
    slug: categoryForm.slug || undefined,
    description: categoryForm.description,
    sort_order: categoryForm.sort_order
  }
  try {
    if (categoryForm.id) {
      await updateCategory(categoryForm.id, payload)
      ElMessage.success('分类已更新')
    } else {
      await createCategory(payload)
      ElMessage.success('分类已新增')
    }
    categoryDialogVisible.value = false
    fetchCategories()
    fetchCategoryStats()
  } catch {
    // 错误提示已由响应拦截器统一处理
  }
}

async function deleteCategory(row) {
  try {
    await ElMessageBox.confirm(`确定删除分类「${row.name}」吗？`, '提示', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消'
    })
  } catch {
    return
  }
  try {
    await removeCategory(row.id)
    ElMessage.success('分类已删除')
    fetchCategories()
    fetchCategoryStats()
  } catch (err) {
    // 该分类下仍有文章时后端返回 409，给出明确提示
    if (err?.code === 409) {
      ElMessage.warning('该分类下仍有文章，无法删除，请先调整文章分类')
    }
  }
}

/* ============ 标签管理 ============ */
const tags = ref([])
const tagsLoading = ref(false)

async function fetchTags() {
  tagsLoading.value = true
  try {
    tags.value = (await listTags()) || []
  } catch {
    tags.value = []
  } finally {
    tagsLoading.value = false
  }
}

const tagDialogVisible = ref(false)
const tagDialogTitle = ref('新增标签')
const tagFormRef = ref()
const tagForm = reactive({
  id: null,
  name: '',
  slug: ''
})
const tagRules = {
  name: [{ required: true, message: '请输入标签名称', trigger: 'blur' }]
}

function openTagDialog(row) {
  if (row) {
    tagDialogTitle.value = '编辑标签'
    tagForm.id = row.id
    tagForm.name = row.name
    tagForm.slug = row.slug
  } else {
    tagDialogTitle.value = '新增标签'
    tagForm.id = null
    tagForm.name = ''
    tagForm.slug = ''
  }
  tagDialogVisible.value = true
}

async function submitTag() {
  if (!tagFormRef.value) return
  try {
    await tagFormRef.value.validate()
  } catch {
    return
  }
  const payload = {
    name: tagForm.name,
    slug: tagForm.slug || undefined
  }
  try {
    if (tagForm.id) {
      await updateTag(tagForm.id, payload)
      ElMessage.success('标签已更新')
    } else {
      await createTag(payload)
      ElMessage.success('标签已新增')
    }
    tagDialogVisible.value = false
    fetchTags()
  } catch {
    // 错误提示已由响应拦截器统一处理
  }
}

async function deleteTag(row) {
  try {
    await ElMessageBox.confirm(`确定删除标签「${row.name}」吗？`, '提示', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消'
    })
  } catch {
    return
  }
  try {
    await removeTag(row.id)
    ElMessage.success('标签已删除')
    fetchTags()
  } catch {
    // 错误提示已由响应拦截器统一处理
  }
}

onMounted(() => {
  fetchOverview()
  fetchCategoryStats()
  fetchCategories()
  fetchTags()
})
</script>

<template>
  <el-card shadow="never" class="admin-card">
    <el-tabs v-model="activeTab">
      <!-- 概览 -->
      <el-tab-pane label="概览" name="overview">
        <div v-loading="overviewLoading">
          <div v-if="overview" class="stat-grid">
            <el-card
              v-for="item in statCards"
              :key="item.label"
              shadow="never"
              class="stat-card"
            >
              <div class="stat-value" :style="{ color: item.color }">{{ item.value }}</div>
              <div class="stat-label">{{ item.label }}</div>
            </el-card>
          </div>
          <el-empty
            v-else-if="!overviewLoading && overviewUnimplemented"
            description="概览统计暂不可用：后端 StatsService 待实现"
          />
          <el-empty v-else-if="!overviewLoading" description="暂无统计数据" />
        </div>

        <h3 class="section-title">分类维度统计</h3>
        <div v-loading="categoryStatsLoading">
          <el-table v-if="categoryStats.length" :data="categoryStats" class="stats-table">
            <el-table-column prop="name" label="分类" width="200" />
            <el-table-column prop="article_count" label="文章数" width="100" />
            <el-table-column label="占比">
              <template #default="{ row }">
                <el-progress
                  :percentage="progressPercent(row.article_count)"
                  :stroke-width="14"
                  :text-inside="true"
                />
              </template>
            </el-table-column>
          </el-table>
          <el-empty
            v-else-if="!categoryStatsLoading && categoryStatsUnimplemented"
            description="分类统计暂不可用：后端 StatsService 待实现"
          />
          <el-empty v-else-if="!categoryStatsLoading" description="暂无分类统计" />
        </div>
      </el-tab-pane>

      <!-- 分类管理 -->
      <el-tab-pane label="分类管理" name="categories">
        <div class="tab-toolbar">
          <el-button type="primary" @click="openCategoryDialog(null)">
            <el-icon><Plus /></el-icon>
            <span>新增分类</span>
          </el-button>
        </div>
        <el-table v-loading="categoriesLoading" :data="categories" empty-text="暂无分类">
          <el-table-column prop="name" label="名称" min-width="140" />
          <el-table-column prop="slug" label="Slug" min-width="140" />
          <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
          <el-table-column prop="sort_order" label="排序" width="80" />
          <el-table-column prop="article_count" label="文章数" width="90" />
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button text size="small" @click="openCategoryDialog(row)">编辑</el-button>
              <el-button text size="small" type="danger" @click="deleteCategory(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <!-- 标签管理 -->
      <el-tab-pane label="标签管理" name="tags">
        <div class="tab-toolbar">
          <el-button type="primary" @click="openTagDialog(null)">
            <el-icon><Plus /></el-icon>
            <span>新增标签</span>
          </el-button>
        </div>
        <el-table v-loading="tagsLoading" :data="tags" empty-text="暂无标签">
          <el-table-column prop="name" label="名称" min-width="160" />
          <el-table-column prop="slug" label="Slug" min-width="160" />
          <el-table-column prop="article_count" label="文章数" width="90" />
          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button text size="small" @click="openTagDialog(row)">编辑</el-button>
              <el-button text size="small" type="danger" @click="deleteTag(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 分类新增/编辑对话框 -->
    <el-dialog v-model="categoryDialogVisible" :title="categoryDialogTitle" width="480px">
      <el-form ref="categoryFormRef" :model="categoryForm" :rules="categoryRules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="categoryForm.name" placeholder="分类名称" />
        </el-form-item>
        <el-form-item label="Slug">
          <el-input v-model="categoryForm.slug" placeholder="留空由名称自动生成" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="categoryForm.description" placeholder="分类描述" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="categoryForm.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCategory">保存</el-button>
      </template>
    </el-dialog>

    <!-- 标签新增/编辑对话框 -->
    <el-dialog v-model="tagDialogVisible" :title="tagDialogTitle" width="440px">
      <el-form ref="tagFormRef" :model="tagForm" :rules="tagRules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="tagForm.name" placeholder="标签名称" />
        </el-form-item>
        <el-form-item label="Slug">
          <el-input v-model="tagForm.slug" placeholder="留空由名称自动生成" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tagDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTag">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<style scoped>
.admin-card {
  border-radius: var(--radius);
}
.stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}
.stat-card {
  border-radius: var(--radius);
  text-align: center;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}
.stat-label {
  margin-top: 6px;
  color: var(--text-muted);
  font-size: 13px;
}
.section-title {
  margin: 8px 0 16px;
  font-size: 16px;
}
.stats-table {
  width: 100%;
}
.tab-toolbar {
  margin-bottom: 16px;
}
</style>
