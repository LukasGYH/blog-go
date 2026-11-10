# 个人博客 · 前端（Vue 3 + Vite）

技术栈：Vue 3（`<script setup>`）+ Vite + Element Plus（全量引入）+ Vue Router 4 + Pinia + Axios + markdown-it + highlight.js。

## 启动方式

```bash
# 1. 安装依赖
npm install

# 2. 本地开发（默认 http://localhost:5173）
npm run dev

# 3. 生产构建（产出 dist/）
npm run build

# 4. 本地预览构建产物
npm run preview
```

> 要求 Node ≥ 18（本项目在 Node v22 下验证通过）。

## 与后端联调

- 开发环境下，前端统一以 `/api/v1/...` 作为请求前缀，Vite 通过代理转发到后端：
  - 代理配置见 `vite.config.js`：`/api` → `http://localhost:8080`（`changeOrigin: true`）。
- 先启动后端：`cd ../backend && make run`（监听 8080），再启动前端 `npm run dev`。
- 认证方式：登录后 token 写入 `localStorage`（key: `blog_token`），请求拦截器自动注入
  `Authorization: Bearer <token>`。
- 响应体统一结构 `{ code, message, data }`，响应拦截器统一剥壳：`code` 为 `200/201` 时直接返回
  `data`，否则弹出错误提示并 reject。

### 错误与降级策略

| 场景 | 前端行为 |
|------|----------|
| HTTP 401 | 清空本地登录态，提示「登录已过期」，跳转 `/login`（带 `redirect`） |
| HTTP 501（后端功能待实现） | 不跳登录、不崩溃，`ElMessage.warning` 提示，页面以空状态兜底 |
| 网络错误 / 后端未启动 | 提示「后端服务未启动，请先运行 make run」 |
| 业务错误（400/403/404/409/422 等） | 弹出后端返回的 `message`，并 reject 给调用方 |

> 后端骨架中 **Article / Tag / Comment / Stats** 模块的 Service 返回 501（尚未实现），
> 首页列表、文章详情、评论、归档、后台统计等均以 `ElEmpty` 空状态优雅降级展示；
> **Auth 与 Category** 模块后端已实现，登录/注册、分类列表可真实联调。

## 页面清单

| 路由 | 页面 | 鉴权 | 说明 |
|------|------|------|------|
| `/` | 首页 | 否 | 文章卡片流 + 搜索/分类/标签/排序筛选 + 分页；右侧栏分类、标签云、热门文章、归档入口。筛选条件与分页同步到 URL query，刷新不丢状态 |
| `/article/:id` | 文章详情 | 否 | Markdown 渲染、作者/分类/标签、浏览量；评论区（发表/回复/删除，未登录提示登录）；作者/管理员可编辑、删除 |
| `/login` | 登录 | 否 | 表单校验，登录成功后跳回 `redirect` 来源页 |
| `/register` | 注册 | 否 | 用户名/邮箱/密码校验（与后端 binding 一致） |
| `/editor` | 新建文章 | 是 | Markdown 编辑 + 实时预览 + 分类单选 + 标签多选 + 草稿/发布 |
| `/editor/:id` | 编辑文章 | 是 | 同上，进入时回填数据 |
| `/me` | 个人中心 | 是 | 我的文章（草稿/已发布筛选、编辑/删除/切换状态）、个人资料、修改密码 |
| `/admin` | 后台 | 管理员 | 概览统计卡片 + 分类维度统计（进度条）、分类管理、标签管理 |
| `/archive` | 归档 | 否 | 按年分组、按月展开统计数量 |
| `*` | 404 | 否 | 兜底页面 |

## 目录结构

```
frontend/
├── index.html
├── package.json
├── vite.config.js
├── .gitignore
└── src/
    ├── main.js               # 挂载 ElementPlus + 全局图标 + router + pinia
    ├── App.vue
    ├── styles/index.css      # 全局样式（浅色主题、卡片阴影、容器最大宽度 1200px）
    ├── utils/format.js       # 时间格式化 / 相对时间 / 摘要截断 / 首字母
    ├── api/                  # request.js（拦截器）+ 各模块接口封装
    ├── stores/user.js        # Pinia：token + user，持久化到 localStorage
    ├── router/index.js       # 路由 + 前置守卫（requiresAuth / requiresAdmin）
    ├── layouts/DefaultLayout.vue
    ├── components/           # 文章卡片、筛选栏、Markdown 渲染、评论列表/单条
    └── views/                # 各页面
```

## 说明

- 不使用任何占位图片服务；文章无封面时使用小尺寸 CSS 渐变色块占位。
- 前端仅做体验层权限隐藏（按钮显示与否），**真实权限以服务端校验为准**。
- markdown-it 使用默认 `html: false`，可防止正文原始 HTML 注入。
