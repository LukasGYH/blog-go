<script setup>
import { computed } from 'vue'

const props = defineProps({
  keyword: { type: String, default: '' },
  categoryId: { type: [Number, String], default: '' },
  tagId: { type: [Number, String], default: '' },
  sort: { type: String, default: 'latest' },
  categories: { type: Array, default: () => [] },
  tags: { type: Array, default: () => [] }
})

const emit = defineEmits([
  'update:keyword',
  'update:categoryId',
  'update:tagId',
  'update:sort',
  'search'
])

// 关键字双向绑定
const keywordValue = computed({
  get: () => props.keyword,
  set: (v) => emit('update:keyword', v)
})

// 分类双向绑定
const categoryValue = computed({
  get: () => props.categoryId,
  set: (v) => emit('update:categoryId', v)
})

// 标签双向绑定
const tagValue = computed({
  get: () => props.tagId,
  set: (v) => emit('update:tagId', v)
})

// 排序双向绑定
const sortValue = computed({
  get: () => props.sort,
  set: (v) => emit('update:sort', v)
})

function onSearch() {
  emit('search')
}
</script>

<template>
  <el-card class="filter-bar" shadow="never" :body-style="{ padding: '16px 20px' }">
    <div class="filter-row">
      <el-input
        v-model="keywordValue"
        class="keyword-input"
        placeholder="搜索文章标题或摘要，回车确认"
        clearable
        @keyup.enter="onSearch"
        @clear="onSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
        <template #append>
          <el-button @click="onSearch">搜索</el-button>
        </template>
      </el-input>

      <el-select
        v-model="categoryValue"
        class="filter-select"
        placeholder="全部分类"
        clearable
        @change="onSearch"
      >
        <el-option
          v-for="c in categories"
          :key="c.id"
          :label="c.article_count != null ? `${c.name}（${c.article_count}）` : c.name"
          :value="c.id"
        />
      </el-select>

      <el-select
        v-model="tagValue"
        class="filter-select"
        placeholder="全部标签"
        clearable
        filterable
        @change="onSearch"
      >
        <el-option
          v-for="t in tags"
          :key="t.id"
          :label="t.article_count != null ? `${t.name}（${t.article_count}）` : t.name"
          :value="t.id"
        />
      </el-select>

      <el-radio-group v-model="sortValue" class="sort-group" @change="onSearch">
        <el-radio-button value="latest">最新</el-radio-button>
        <el-radio-button value="hot">最热</el-radio-button>
      </el-radio-group>
    </div>
  </el-card>
</template>

<style scoped>
.filter-bar {
  border-radius: var(--radius);
  margin-bottom: 16px;
}
.filter-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
.keyword-input {
  flex: 1;
  min-width: 220px;
}
.filter-select {
  width: 180px;
}
.sort-group {
  margin-left: auto;
}
@media (max-width: 768px) {
  .filter-select {
    width: 100%;
  }
  .sort-group {
    margin-left: 0;
  }
}
</style>
