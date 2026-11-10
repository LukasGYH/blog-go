<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createArticle, updateArticle, getArticle } from '@/api/article'
import { listCategories } from '@/api/category'
import { listTags } from '@/api/tag'
import MarkdownView from '@/components/MarkdownView.vue'

const route = useRoute()
const router = useRouter()

const formRef = ref()
const saving = ref(false)
const loading = ref(false)

const categories = ref([])
const tags = ref([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  summary: '',
  content: '',
  cover_image: '',
  category_id: null,
  tag_ids: [],
  status: 'draft'
})

const rules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { max: 200, message: '标题不能超过 200 个字符', trigger: 'blur' }
  ],
  content: [{ required: true, message: '请输入正文内容', trigger: 'blur' }]
}

const previewContent = computed(() => form.content || '（正文预览区）')

async function loadOptions() {
  try {
    categories.value = (await listCategories()) || []
  } catch {
    categories.value = []
  }
  try {
    tags.value = (await listTags()) || []
  } catch {
    tags.value = []
  }
}

// 编辑模式：回填数据（注意后端返回的 category / tags / status 需要正确映射）
async function loadArticle() {
  if (!isEdit.value) return
  loading.value = true
  try {
    const data = await getArticle(route.params.id)
    form.title = data.title || ''
    form.summary = data.summary || ''
    form.content = data.content || ''
    form.cover_image = data.cover_image || ''
    form.category_id = data.category?.id ?? null
    form.tag_ids = (data.tags || []).map((t) => t.id)
    form.status = data.status || 'draft'
  } catch {
    // 501 或 404 等：错误提示已由响应拦截器处理，页面保持空表单
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  const payload = {
    title: form.title,
    summary: form.summary,
    content: form.content,
    cover_image: form.cover_image,
    category_id: form.category_id || 0,
    tag_ids: form.tag_ids,
    status: form.status
  }

  saving.value = true
  try {
    const data = isEdit.value
      ? await updateArticle(route.params.id, payload)
      : await createArticle(payload)
    ElMessage.success('保存成功')
    const newId = data?.id || route.params.id
    if (newId) router.push(`/article/${newId}`)
    else router.push('/me')
  } catch {
    // 错误提示已由响应拦截器统一处理（含 501 待实现提示）
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await loadOptions()
  await loadArticle()
})
</script>

<template>
  <div v-loading="loading" class="editor-page">
    <div class="editor-header">
      <h2 class="editor-title">{{ isEdit ? '编辑文章' : '写文章' }}</h2>
      <div class="editor-header-actions">
        <el-radio-group v-model="form.status">
          <el-radio-button value="draft">草稿</el-radio-button>
          <el-radio-button value="published">发布</el-radio-button>
        </el-radio-group>
        <el-button type="primary" :loading="saving" @click="handleSave">
          {{ form.status === 'published' ? '保存并发布' : '保存草稿' }}
        </el-button>
      </div>
    </div>

    <div class="editor-body">
      <!-- 编辑区 -->
      <div class="editor-left">
        <el-card shadow="never" class="editor-card">
          <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
            <el-form-item label="标题" prop="title">
              <el-input v-model="form.title" placeholder="请输入文章标题" maxlength="200" show-word-limit />
            </el-form-item>

            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="分类">
                  <el-select
                    v-model="form.category_id"
                    placeholder="选择分类（可选）"
                    clearable
                    class="full-width"
                  >
                    <el-option
                      v-for="c in categories"
                      :key="c.id"
                      :label="c.name"
                      :value="c.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="封面图 URL">
                  <el-input v-model="form.cover_image" placeholder="图片 URL（可选，仅存字符串）" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="标签">
              <el-select
                v-model="form.tag_ids"
                multiple
                filterable
                collapse-tags
                collapse-tags-tooltip
                placeholder="选择标签（可输入过滤，可多选）"
                class="full-width"
              >
                <el-option
                  v-for="t in tags"
                  :key="t.id"
                  :label="t.name"
                  :value="t.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="摘要">
              <el-input
                v-model="form.summary"
                type="textarea"
                :rows="2"
                maxlength="500"
                placeholder="留空则由后端自动从正文截取前 120 字生成"
              />
            </el-form-item>

            <el-form-item label="正文（Markdown）" prop="content">
              <el-input
                v-model="form.content"
                type="textarea"
                :rows="20"
                placeholder="支持 Markdown 语法，右侧实时预览"
              />
            </el-form-item>
          </el-form>
        </el-card>
      </div>

      <!-- 实时预览 -->
      <div class="editor-right">
        <el-card shadow="never" class="editor-card preview-card">
          <template #header>
            <span class="card-title" style="margin: 0">实时预览</span>
          </template>
          <h3 class="preview-title">{{ form.title || '未命名文章' }}</h3>
          <MarkdownView :content="previewContent" />
        </el-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}
.editor-title {
  margin: 0;
  font-size: 20px;
}
.editor-header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
.editor-body {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}
.editor-left {
  flex: 1;
  min-width: 0;
}
.editor-right {
  width: 460px;
  flex-shrink: 0;
  position: sticky;
  top: 84px;
}
.editor-card {
  border-radius: var(--radius);
}
.full-width {
  width: 100%;
}
.preview-card {
  max-height: calc(100vh - 120px);
  overflow: auto;
}
.preview-title {
  margin: 0 0 12px;
  font-size: 20px;
}
@media (max-width: 1100px) {
  .editor-body {
    flex-direction: column;
  }
  .editor-right {
    width: 100%;
    position: static;
  }
}
</style>
