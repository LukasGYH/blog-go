<script setup>
import { computed } from 'vue'
import { formatDateTime, truncate, initialOf } from '@/utils/format'

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

// 摘要优先用后端 summary，缺省时从标题兜底
const summary = computed(() => {
  const a = props.article
  return truncate(a.summary || '', 140) || '（暂无摘要）'
})

const authorName = computed(
  () => props.article.author?.nickname || props.article.author?.username || '匿名'
)

const avatarLetter = computed(() => initialOf(authorName.value))

const tags = computed(() => props.article.tags || [])
</script>

<template>
  <el-card class="article-card" shadow="never" :body-style="{ padding: '18px 20px' }">
    <router-link :to="`/article/${article.id}`" class="card-link">
      <div class="card-body">
        <div class="card-main">
          <h3 class="article-title">{{ article.title }}</h3>
          <p class="article-summary">{{ summary }}</p>

          <div class="article-meta">
            <span class="meta-author">
              <span v-if="article.author?.avatar" class="mini-avatar">
                <img :src="article.author.avatar" alt="avatar" />
              </span>
              <span v-else class="avatar-letter mini-avatar-letter">
                {{ avatarLetter }}
              </span>
              {{ authorName }}
            </span>

            <el-tag v-if="article.category" size="small" type="info" effect="plain">
              {{ article.category.name }}
            </el-tag>

            <el-tag
              v-for="tag in tags.slice(0, 3)"
              :key="tag.id"
              size="small"
              effect="plain"
            >
              {{ tag.name }}
            </el-tag>

            <span class="meta-item">
              <el-icon><View /></el-icon>{{ article.view_count ?? 0 }}
            </span>
            <span class="meta-item">
              <el-icon><ChatDotRound /></el-icon>{{ article.comment_count ?? 0 }}
            </span>
            <span class="meta-item">
              <el-icon><Clock /></el-icon>
              {{ formatDateTime(article.published_at || article.created_at) }}
            </span>

            <el-tag v-if="article.status === 'draft'" size="small" type="warning">
              草稿
            </el-tag>
          </div>
        </div>

        <div class="card-cover">
          <img v-if="article.cover_image" :src="article.cover_image" alt="cover" />
          <div v-else class="cover-placeholder">BLOG</div>
        </div>
      </div>
    </router-link>
  </el-card>
</template>

<style scoped>
.article-card {
  border-radius: var(--radius);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.article-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.08);
}
.card-link {
  display: block;
}
.card-body {
  display: flex;
  gap: 16px;
}
.card-main {
  flex: 1;
  min-width: 0;
}
.article-title {
  margin: 0 0 10px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-main);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.article-card:hover .article-title {
  color: var(--primary);
}
.article-summary {
  margin: 0 0 14px;
  color: var(--text-sub);
  line-height: 1.7;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.article-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  color: var(--text-muted);
  font-size: 13px;
}
.meta-author {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--text-sub);
}
.mini-avatar {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  overflow: hidden;
}
.mini-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.mini-avatar-letter {
  width: 22px;
  height: 22px;
  font-size: 12px;
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.card-cover {
  width: 160px;
  height: 100px;
  border-radius: 6px;
  overflow: hidden;
  flex-shrink: 0;
}
.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
@media (max-width: 640px) {
  .card-cover {
    display: none;
  }
}
</style>
