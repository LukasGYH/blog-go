<script setup>
import { ref, computed, onMounted } from 'vue'
import { getArchive } from '@/api/article'
import { useRouter } from 'vue-router'

const router = useRouter()

const archive = ref([])
const loading = ref(false)
const unimplemented = ref(false)
// 展开的月份 key
const expandedMonths = ref([])

// 按年分组：[{ year, months: [{ month, count }] }]
const grouped = computed(() => {
  const map = new Map()
  for (const item of archive.value) {
    if (!map.has(item.year)) map.set(item.year, [])
    map.get(item.year).push({ month: item.month, count: item.count })
  }
  return Array.from(map.entries())
    .map(([year, months]) => ({
      year,
      months: months.sort((a, b) => b.month - a.month),
      total: months.reduce((sum, m) => sum + (m.count || 0), 0)
    }))
    .sort((a, b) => b.year - a.year)
})

async function fetchArchive() {
  loading.value = true
  try {
    archive.value = (await getArchive()) || []
    unimplemented.value = false
  } catch (err) {
    archive.value = []
    if (err?.code === 501) unimplemented.value = true
  } finally {
    loading.value = false
  }
}

function monthKey(year, month) {
  return `${year}-${month}`
}

// 展开月份面板
function toggleMonth(year, month) {
  const key = monthKey(year, month)
  if (expandedMonths.value.includes(key)) {
    expandedMonths.value = expandedMonths.value.filter((k) => k !== key)
  } else {
    expandedMonths.value = [...expandedMonths.value, key]
  }
}

function isExpanded(year, month) {
  return expandedMonths.value.includes(monthKey(year, month))
}

function goHome() {
  router.push('/')
}

onMounted(fetchArchive)
</script>

<template>
  <div class="archive-page">
    <el-card shadow="never" class="archive-card">
      <h2 class="archive-heading">
        <el-icon><Calendar /></el-icon>
        <span>文章归档</span>
      </h2>
      <p class="archive-desc">按年月浏览全部已发布文章</p>

      <div v-loading="loading">
        <el-collapse v-if="grouped.length" class="year-collapse">
          <el-collapse-item v-for="group in grouped" :key="group.year" :name="group.year">
            <template #title>
              <span class="year-title">{{ group.year }} 年</span>
              <el-tag size="small" type="info" effect="plain" class="year-count">
                {{ group.total }} 篇
              </el-tag>
            </template>

            <ul class="month-list">
              <li
                v-for="m in group.months"
                :key="m.month"
                class="month-item"
                :class="{ 'month-item--open': isExpanded(group.year, m.month) }"
              >
                <div class="month-head" @click="toggleMonth(group.year, m.month)">
                  <span class="month-name">{{ m.month }} 月</span>
                  <span class="month-count">{{ m.count }} 篇</span>
                  <el-icon class="month-arrow">
                    <ArrowRight v-if="!isExpanded(group.year, m.month)" />
                    <ArrowDown v-else />
                  </el-icon>
                </div>

                <transition name="el-fade-in">
                  <div v-if="isExpanded(group.year, m.month)" class="month-body">
                    <!-- 说明：当前后端契约（GET /articles）不支持按月筛选，
                         为避免展示未经筛选的错误数据，此处仅展示数量与提示 -->
                    <el-alert
                      :title="`该月共 ${m.count} 篇已发布文章`"
                      type="info"
                      :closable="false"
                      show-icon
                    >
                      <template #default>
                        <div class="month-tip">
                          后端暂未提供「按月份筛选文章列表」的接口，无法直接展开明细。
                          可返回首页按关键字 / 分类 / 标签检索文章。
                        </div>
                        <el-button size="small" type="primary" plain @click="goHome">
                          去首页浏览
                        </el-button>
                      </template>
                    </el-alert>
                  </div>
                </transition>
              </li>
            </ul>
          </el-collapse-item>
        </el-collapse>

        <el-empty
          v-else-if="!loading && unimplemented"
          description="归档统计暂不可用：后端 ArticleService 待实现"
        />
        <el-empty v-else-if="!loading" description="暂无归档数据" />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.archive-page {
  max-width: 900px;
  margin: 0 auto;
}
.archive-card {
  border-radius: var(--radius);
}
.archive-heading {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 6px;
  font-size: 22px;
}
.archive-desc {
  margin: 0 0 20px;
  color: var(--text-muted);
  font-size: 13px;
}
.year-title {
  font-size: 16px;
  font-weight: 600;
}
.year-count {
  margin-left: 12px;
}
.month-list {
  list-style: none;
  margin: 0;
  padding: 0 0 0 8px;
}
.month-item {
  border-bottom: 1px dashed var(--border-light);
}
.month-item:last-child {
  border-bottom: none;
}
.month-head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 8px;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.2s;
}
.month-head:hover {
  background: #f2f6fc;
}
.month-name {
  font-weight: 500;
  color: var(--text-main);
  min-width: 64px;
}
.month-count {
  color: var(--text-muted);
  font-size: 13px;
}
.month-arrow {
  margin-left: auto;
  color: var(--text-muted);
}
.month-body {
  padding: 4px 8px 16px;
}
.month-tip {
  margin: 6px 0 10px;
  line-height: 1.7;
}
</style>
