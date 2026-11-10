# blog-api 后端骨架

个人博客系统的后端骨架，技术栈为 **Gin + GORM + SQLite + JWT**，采用
`handler → service → repository → entity` 的分层结构，统一响应体与统一错误处理。

API 契约以 `../docs/01-需求分析.md` 第 6 章（冻结版）为准。

## 环境要求

- Go 1.22 及以上（本项目在 go 1.27 环境验证通过）
- CGO 可用（SQLite 使用 CGO 驱动 `gorm.io/driver/sqlite`）
- 国内网络建议设置代理：`go env -w GOPROXY=https://goproxy.cn,direct`

## 快速开始

```bash
# 1. 准备配置（首次）
cp .env.example .env

# 2. 整理依赖
go mod tidy

# 3. 启动服务（默认监听 0.0.0.0:8080，数据库文件 blog.db）
make run
# 等价于：go run .
```

启动成功后会打印监听地址与数据库文件路径：

```
服务监听地址: http://0.0.0.0:8080
数据库文件路径: blog.db
```

首次启动会自动执行 `AutoMigrate` 建表，并**幂等**灌入种子数据（已存在则跳过）。

### 常用命令（Makefile）

| 命令 | 说明 |
| --- | --- |
| `make run` | 启动服务 |
| `make build` | 编译到 `bin/blog-api` |
| `make tidy` | 整理依赖 |
| `make test` | 运行测试 |
| `make reset-db` | 删除 `blog.db`，下次启动自动重建并重灌种子 |

### 配置项（.env）

| 变量 | 说明 | 默认 |
| --- | --- | --- |
| `SERVER_HOST` | 监听地址 | `0.0.0.0` |
| `SERVER_PORT` | 监听端口 | `8080` |
| `SERVER_MODE` | Gin 运行模式（debug/release） | `debug` |
| `DB_PATH` | SQLite 文件路径 | `blog.db` |
| `JWT_SECRET` | JWT 签名密钥（生产环境务必更换） | `change-me-in-production` |
| `JWT_EXPIRE` | Token 有效期（如 `24h`） | `24h` |

## 默认账号（种子数据）

| 用户名 | 密码 | 角色 | 邮箱 |
| --- | --- | --- | --- |
| `admin` | `admin123` | admin（管理员） | admin@example.com |
| `lukas` | `123456` | user（普通用户） | lukas@example.com |

种子数据还包含：3 个分类、5 个标签、4 篇文章（3 篇已发布 + 1 篇草稿，草稿作者为 lukas）。

## 模块实现状态

### 已实现（完整可用）

- **Auth 认证**：注册（唯一性预检 + bcrypt）、登录（失败统一返回 401）、
  获取/更新个人资料、修改密码
- **Category 分类**：公开列表（含 `article_count`）、创建、更新、删除
  （分类下存在文章时返回 409）

这两个模块同时作为 repo + service 的**参考样板**，可对照实现其余模块。

### 待实现（当前返回 501）

以下方法为**空实现**，调用时统一返回 `code=501`、`message=功能待实现：<模块>.<方法>`。
每个方法上方均留有中文 TODO 注释，说明要实现的功能与涉及的 GORM/Gin 知识点
（编号对应需求文档第 11 章练习任务清单）。

| 模块 | 待实现方法 |
| --- | --- |
| Handler 层 | `ArticleHandler`(8)、`TagHandler`(4)、`CommentHandler`(3)、`StatsHandler`(2)，共 **17 个** |
| Tag 标签 | `TagRepository.List/GetByID/Create/Update/Delete`、`TagService.List/Create/Update/Delete` |
| Article 文章 | `ArticleRepository.List/GetByID/GetBySlug/Create/Update/Delete/UpdateStatus/IncrViewCount/Archive/Hot/CountByCategory` |
| Article 文章 | `ArticleService.List/GetByID/Create/Update/Delete/ChangeStatus/Archive/Hot` |
| Comment 评论 | `CommentRepository.ListTopLevel/ListReplies/GetByID/Create/Delete/DeleteReplies`、`CommentService.ListByArticle/Create/Delete` |
| Stats 统计 | `StatsRepository.Overview/CountByCategory`、`StatsService.Overview/CategoryStats` |

> 说明：文章列表等接口返回 501 属于**预期行为**，是留给练习者的实现位。

## 接口速览

- 健康检查：`GET /health` → 纯文本 `ok`（顶层，不在 `/api/v1` 下）
- 认证：`/api/v1/auth/register|login|me|password`
- 分类：`/api/v1/categories`
- 标签：`/api/v1/tags`
- 文章：`/api/v1/articles`、`/api/v1/articles/archive`、`/api/v1/articles/hot`
- 评论：`/api/v1/articles/:id/comments`、`/api/v1/comments/:id`
- 统计：`/api/v1/stats/overview`、`/api/v1/stats/categories`（仅管理员）

鉴权规则：文章读接口为可选鉴权（游客可见已发布文章），写操作需登录；
分类/标签写操作与统计接口仅管理员可访问。

## 查看 SQLite 数据

```bash
# 使用 sqlite3 命令行
sqlite3 blog.db ".tables"
sqlite3 blog.db "SELECT id, username, role FROM users;"
sqlite3 blog.db "SELECT id, title, status FROM articles;"
```

## 目录结构

```
backend/
├── main.go              # 入口：加载配置 → 连库 → AutoMigrate → 种子 → 路由 → 启动
├── config/              # 配置读取（.env）
├── constant/            # 业务常量（角色、状态、排序，自定义字符串类型）
├── utils/               # 统一响应、错误、JWT、密码、分页、slug
├── models/
│   ├── entity/          # 数据实体（GORM）
│   ├── dto/             # 请求参数（带 binding 校验）
│   └── vo/              # 响应视图对象（JSON 字段与契约一致）
├── repositories/        # 数据访问层（db.go 负责建库与种子）
├── services/            # 业务逻辑层
├── handlers/            # HTTP 处理层（禁止出现 GORM 代码）
├── middleware/          # 日志、跨域、恢复、鉴权
└── router/              # 路由装配
```
