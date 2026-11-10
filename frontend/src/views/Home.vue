<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { listArticles, getHotArticles } from '@/api/article'
import { listCategories } from '@/api/category'
import { listTags } from '@/api/tag'
import ArticleCard from '@/components/ArticleCard.vue'
import ArticleFilterBar from '@/components/ArticleFilterBar.vue'
import { formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()

// 筛选条件（以 URL query 为单一数据源）
const keyword = ref('')
const categoryId = ref('')
const tagId = ref('')
const sort = ref('latest')
const page = ref(1)
const pageSize = 10

const articles = ref([])
const total = ref(0)
const loading = ref(false)
// 后端 ArticleService 未实现（501）时标记
const unimplemented = ref(false)

// 侧栏数据
const categories = ref([])
const tags = ref([])
const hotArticles = ref([])
const categoriesReady = ref(true)

function buildQuery(targetPage) {
  const q = {}
  if (keyword.value) q.keyword = keyword.value
  if (categoryId.value) q.category_id = categoryId.value
  if (tagId.value) q.tag_id = tagId.value
  if (sort.value && sort.value !== 'latest') q.sort = sort.value
  if (targetPage > 1) q.page = targetPage
  return q
}

// 从 URL query 初始化本地筛选条件（刷新不丢状态）
function syncFromQuery() {
  keyword.value = route.query.keyword || ''
  categoryId.value = route.query.category_id || ''
  tagId.value = route.query.tag_id || ''
  sort.value = route.query.sort || 'latest'
  page.value = Number(route.query.page) || 1
}

function applyFilter() {
  router.push({ path: '/', query: buildQuery(1) })
}

function handlePageChange(p) {
  router.push({ path: '/', query: buildQuery(p) })
}

async function fetchArticles() {
  loading.value = true
  try {
    const data = await listArticles({
      page: page.value,
      page_size: pageSize,
      keyword: keyword.value || undefined,
      category_id: categoryId.value || undefined,
      tag_id: tagId.value || undefined,
      sort: sort.value || undefined
    })
    articles.value = data?.list || []
    total.value = data?.total || 0
    unimplemented.value = false
  } catch (err) {
    articles.value = []
    total.value = 0
    if (err?.code === 501) unimplemented.value = true
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    categories.value = (await listCategories()) || []
    categoriesReady.value = true
  } catch {
    categories.value = []
    categoriesReady.value = false
  }
}

async function fetchTags() {
  try {
    tags.value = (await listTags()) || []
  } catch {
    tags.value = []
  }
}

async function fetchHot() {
  try {
    hotArticles.value = (await getHotArticles(5)) || []
  } catch {
    hotArticles.value = []
  }
}

// 侧栏点击筛选
function pickCategory(id) {
  categoryId.value = categoryId.value === id ? '' : id
  applyFilter()
}
function pickTag(id) {
  tagId.value = tagId.value === id ? '' : id
  applyFilter()
}

watch(
  () => route.query,
  () => {
    syncFromQuery()
    fetchArticles()
  },
  { immediate: true }
)

// 侧栏数据只加载一次
fetchCategories()
fetchTags()
fetchHot()
</script>

<template>
  <div class="home">
    <div class="home-main">
      <ArticleFilterBar
        v-model:keyword="keyword"
        v-model:category-id="categoryId"
        v-model:tag-id="tagId"
        v-model:sort="sort"
        :categories="categories"
        :tags="tags"
        @search="applyFilter"
      />

      <div v-loading="loading" class="article-list">
        <template v-if="articles.length">
          <ArticleCard
            v-for="article in articles"
            :key="article.id"
            :article="article"
            class="article-list-item"
          />
        </template>
        <el-card v-else-if="!loading" shadow="never" class="empty-card">
          <el-empty
            :description="
              unimplemented
                ? '文章列表暂不可用：后端 ArticleService 待实现'
                : '没有找到符合条件的文章'
            "
          />
        </el-card>
      </div>

      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next, total"
          background
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <aside class="home-aside">
      <!-- 分类 -->
      <el-card class="aside-card" shadow="never">
        <h3 class="card-title">分类</h3>
        <ul v-if="categories.length" class="cat-list">
          <li
            v-for="c in categories"
            :key="c.id"
            class="cat-item"
            :class="{ 'cat-item--active': String(categoryId) === String(c.id) }"
            @click="pickCategory(c.id)"
          >
            <span class="cat-name">{{ c.name }}</span>
            <el-tag size="small" type="info" effect="plain" round>
              {{ c.article_count ?? 0 }}
            </el-tag>
          </li>
        </ul>
        <el-empty v-else-if="categoriesReady" description="暂无分类" :image-size="60" />
        <el-empty v-else description="分类后端待实现" :image-size="60" />
      </el-card>

      <!-- 标签云 -->
      <el-card class="aside-card" shadow="never">
        <h3 class="card-title">标签</h3>
        <div v-if="tags.length" class="tag-cloud">
          <el-tag
            v-for="t in tags"
            :key="t.id"
            class="tag-cloud-item"
            :type="String(tagId) === String(t.id) ? 'primary' : 'info'"
            :effect="String(tagId) === String(t.id) ? 'dark' : 'plain'"
            @click="pickTag(t.id)"
          >
            {{ t.name }}
          </el-tag>
        </div>
        <el-empty v-else description="标签后端待实现" :image-size="60" />
      </el-card>

      <!-- 热门文章 -->
      <el-card class="aside-card" shadow="never">
        <h3 class="card-title">热门文章</h3>
        <ul v-if="hotArticles.length" class="hot-list">
          <li v-for="(a, index) in hotArticles" :key="a.id" class="hot-item">
            <span class="hot-index" :class="{ 'hot-index--top': index < 3 }">
              {{ index + 1 }}
            </span>
            <router-link :to="`/article/${a.id}`" class="hot-title">
              {{ a.title }}
            </router-link>
            <span class="hot-view">
              <el-icon><View /></el-icon>{{ a.view_count ?? 0 }}
            </span>
          </li>
        </ul>
        <el-empty v-else description="热门文章后端待实现" :image-size="60" />
      </el-card>

      <!-- 归档入口 -->
      <el-card class="aside-card archive-entry" shadow="never" @click="router.push('/archive')">
        <div class="archive-entry-inner">
          <el-icon class="archive-icon"><Calendar /></el-icon>
          <div>
            <div class="archive-title">文章归档</div>
            <div class="archive-sub">按年月浏览全部文章</div>
          </div>
          <el-icon class="archive-arrow"><ArrowRight /></el-icon>
        </div>
      </el-card>
    </aside>
  </div>
</template>

<style scoped>
.home {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}
.home-main {
  flex: 1;
  min-width: 0;
}
.article-list-item {
  margin-bottom: 16px;
}
.empty-card {
  border-radius: var(--radius);
}
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.home-aside {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.aside-card {
  border-radius: var(--radius);
}

.cat-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.cat-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  color: var(--text-sub);
  transition: background 0.2s;
}
.cat-item:hover {
  background: #f2f6fc;
  color: var(--primary);
}
.cat-item--active {
  background: #ecf5ff;
  color: var(--primary);
  font-weight: 600;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag-cloud-item {
  cursor: pointer;
}

.hot-list {
  list-style: none;
  margin: 0;
  padding: 0;
}
.hot-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 0;
  border-bottom: 1px dashed var(--border-light);
}
.hot-item:last-child {
  border-bottom: none;
}
.hot-index {
  width: 18px;
  height: 18px;
  line-height: 18px;
  text-align: center;
  border-radius: 4px;
  font-size: 12px;
  background: #f0f2f5;
  color: var(--text-muted);
  flex-shrink: 0;
}
.hot-index--top {
  background: var(--primary);
  color: #fff;
}
.hot-title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-sub);
}
.hot-title:hover {
  color: var(--primary);
}
.hot-view {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  color: var(--text-muted);
  font-size: 12px;
}

.archive-entry {
  cursor: pointer;
}
.archive-entry-inner {
  display: flex;
  align-items: center;
  gap: 12px;
}
.archive-icon {
  font-size: 24px;
  color: var(--primary);
}
.archive-title {
  font-weight: 600;
  color: var(--text-main);
}
.archive-sub {
  font-size: 12px;
  color: var(--text-muted);
}
.archive-arrow {
  margin-left: auto;
  color: var(--text-muted);
}

@media (max-width: 900px) {
  .home {
    flex-direction: column;
  }
  .home-aside {
    width: 100%;
  }
}
</style>
