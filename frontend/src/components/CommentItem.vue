<script setup>
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { fromNow, initialOf } from '@/utils/format'

const props = defineProps({
  comment: { type: Object, required: true },
  articleAuthorId: { type: [Number, String], default: null }
})

const emit = defineEmits(['reply', 'remove'])

const userStore = useUserStore()

const replies = computed(() => props.comment.replies || [])

function nameOf(user) {
  return user?.nickname || user?.username || '匿名'
}

// 删除按钮仅对「评论作者 / 文章作者 / 管理员」显示（前端仅做体验层隐藏，权限以服务端为准）
function canDelete(comment) {
  if (!userStore.isLogin) return false
  const uid = userStore.user?.id
  if (userStore.isAdmin) return true
  if (comment.user?.id != null && comment.user.id === uid) return true
  if (props.articleAuthorId != null && props.articleAuthorId === uid) return true
  return false
}
</script>

<template>
  <div class="comment-item">
    <div class="comment-head">
      <span v-if="comment.user?.avatar" class="c-avatar">
        <img :src="comment.user.avatar" alt="avatar" />
      </span>
      <span v-else class="avatar-letter c-avatar-letter">
        {{ initialOf(nameOf(comment.user)) }}
      </span>
      <div class="c-meta">
        <span class="c-name">{{ nameOf(comment.user) }}</span>
        <span class="c-time">{{ fromNow(comment.created_at) }}</span>
      </div>
    </div>

    <div class="comment-content">{{ comment.content }}</div>

    <div class="comment-actions">
      <el-button v-if="userStore.isLogin" text size="small" @click="emit('reply', comment)">
        回复
      </el-button>
      <el-button
        v-if="canDelete(comment)"
        text
        size="small"
        type="danger"
        @click="emit('remove', comment)"
      >
        删除
      </el-button>
    </div>

    <!-- 内嵌回复（两级） -->
    <div v-if="replies.length" class="comment-replies">
      <div v-for="reply in replies" :key="reply.id" class="reply-item">
        <div class="comment-head">
          <span v-if="reply.user?.avatar" class="c-avatar">
            <img :src="reply.user.avatar" alt="avatar" />
          </span>
          <span v-else class="avatar-letter c-avatar-letter">
            {{ initialOf(nameOf(reply.user)) }}
          </span>
          <div class="c-meta">
            <span class="c-name">{{ nameOf(reply.user) }}</span>
            <span class="c-time">{{ fromNow(reply.created_at) }}</span>
          </div>
        </div>
        <div class="comment-content">{{ reply.content }}</div>
        <div class="comment-actions">
          <el-button v-if="userStore.isLogin" text size="small" @click="emit('reply', reply)">
            回复
          </el-button>
          <el-button
            v-if="canDelete(reply)"
            text
            size="small"
            type="danger"
            @click="emit('remove', reply)"
          >
            删除
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.comment-item {
  padding: 16px 0;
  border-bottom: 1px solid var(--border-light);
}
.comment-item:last-child {
  border-bottom: none;
}
.comment-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}
.c-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}
.c-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.c-avatar-letter {
  width: 36px;
  height: 36px;
  font-size: 15px;
}
.c-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}
.c-name {
  font-weight: 600;
  color: var(--text-main);
}
.c-time {
  font-size: 12px;
  color: var(--text-muted);
}
.comment-content {
  color: var(--text-sub);
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
  padding-left: 46px;
}
.comment-actions {
  padding-left: 46px;
  margin-top: 4px;
}
.comment-replies {
  margin: 12px 0 0 46px;
  padding-left: 16px;
  border-left: 2px solid var(--border-light);
}
.reply-item {
  padding: 10px 0;
}
.reply-item + .reply-item {
  border-top: 1px dashed var(--border-light);
}
</style>
