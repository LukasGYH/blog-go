<script setup>
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const props = defineProps({
  content: { type: String, default: '' }
})

// markdown-it 默认 html:false，可有效防止正文中的原始 HTML 注入（XSS 基本处理）
const md = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: true,
  highlight(str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, {
          language: lang,
          ignoreIllegals: true
        }).value}</code></pre>`
      } catch {
        // 降级到下面的转义分支
      }
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  }
})

const rendered = computed(() => md.render(props.content || ''))
</script>

<template>
  <div class="markdown-body" v-html="rendered"></div>
</template>
