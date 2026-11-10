# 个人博客系统（blog-app）

第二部分作业 · Web 应用实战。后端用 Gin + GORM 手写，前端 Vue3 + Element Plus，数据库 SQLite。

**需求分析文档：[docs/01-需求分析.md](file:///Volumes/out_share/cc-space/blockchain/go/lessons/blog-app/docs/01-需求分析.md)** —— 数据模型、API 契约、页面设计、验收标准、练习清单都在里面，是全项目的唯一权威来源。

## 目录结构

```
blog-app/
├── docs/01-需求分析.md      # 需求 / 契约 / 验收 / 练习清单
├── backend/                 # Go 后端（module: blog-api）
└── frontend/                # Vue3 + Vite 前端
```

## 环境要求

| 组件 | 版本 | 本机状态 |
|------|------|----------|
| Go | ≥ 1.22（本机 1.27.3） | 已就绪 |
| GOPROXY | `https://goproxy.cn,direct` | 已配置 |
| Node / npm | ≥ 18（本机 22.22.1 / 10.9.4） | 已就绪 |
| SQLite | `gorm.io/driver/sqlite`（CGO），本机自带 `sqlite3` CLI | 已就绪 |

## 启动

**1. 后端**（首次会自动建表 + 灌种子数据）

```bash
cd backend
make tidy      # 首次拉依赖，会编译 CGO 版 SQLite 驱动，耐心等待
make run       # 监听 http://localhost:8080
```

**2. 前端**

```bash
cd frontend
npm install
npm run dev    # http://localhost:5173 （/api 已代理到 8080）
```

浏览器打开 http://localhost:5173 即可。

**默认账号**

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | `admin` | `admin123` |
| 普通用户 | `lukas` | `123456` |

**其他命令**

```bash
cd backend
make build      # 编译到 bin/blog-api
make test       # 跑单测
make reset-db   # 删除 blog.db，下次启动重建 + 重新灌种子
```

```bash
# 直接看数据库
sqlite3 backend/blog.db ".tables"
sqlite3 backend/blog.db "select id,title,status,view_count from articles;"
```

## 已经能跑通什么 / 还需要你写什么

| 模块 | 后端状态 | 说明 |
|------|----------|------|
| 基础设施（配置 / 中间件 / 路由 / 工具 / 模型 / 常量） | ✅ 已实现 | 直接可用 |
| **Auth**（注册 / 登录 / JWT / 资料 / 改密） | ✅ 已实现 | handler + service + repo 三层全通，认证类**参考样板** |
| **Category**（分类 CRUD + 删除关联校验） | ✅ 已实现 | handler + service + repo 三层全通，标准 CRUD + 聚合**参考样板** |
| Handlers（Article / Tag / Comment / Stats，共 17 个方法） | ⏳ 待你实现 | 空实现返回 501，含 TODO |
| **Article**（文章 repo 11 个 + service 8 个方法） | ⏳ 待你实现 | 返回 501 + TODO 注释 |
| **Tag**（标签 repo + service） | ⏳ 待你实现 | 返回 501 + TODO 注释 |
| **Comment**（评论 repo + service） | ⏳ 待你实现 | 返回 501 + TODO 注释 |
| **Stats**（统计 repo + service） | ⏳ 待你实现 | 返回 501 + TODO 注释 |
| 前端全部页面 | ✅ 已实现 | 遇 501 会降级提示，不会白屏 |

> 练习空间说明：除 **Auth / Category** 两个样板外，**其余模块的 Handler、Service、Repository 三层都是空的**（方法签名已固化 = 你的接口契约），方法体统一返回 501。清点一下待实现方法：

```bash
cd backend
grep -rn "NotImplemented\|TODO" handlers services repositories | wc -l
```

待实现的每个方法上方都有中文 TODO 注释，并标注了它在需求文档第 11 章的练习编号，例如：

```
handlers/article_handler.go
  // GetByID 文章详情：可选鉴权 —— 未发布文章仅作者/管理员可见
  //   - 参考文档第 11 章：11.3 #2
services/article_service.go
  // Create 新建文章：参数补全（slug 唯一化、摘要截取）+ 事务编排
  //   - 参考文档第 11 章：11.2 #12 Create（参数补全、事务编排、错误包装 %w）
```

## 推荐的练习顺序

1. **Tag 模块** —— 最简单的 CRUD，热身（唯一约束冲突识别、中间表清理）
2. **Article 模块 GetByID + IncrViewCount** —— 嵌套 Preload、`gorm.Expr` 原子自增
3. **Article 模块 List** —— 多条件动态查询 + 分页 + 排序 + 字段裁剪 + `comment_count` 统计
4. **Article 模块 Create / Update** —— 事务内重设多对多关联、越权校验、首次发布写 `published_at`
5. **Comment 模块** —— 两级评论树组装（map 索引 + 一次遍历）
6. **Stats + Archive** —— 聚合查询（`GROUP BY` + `Scan` 到自定义 VO）

每完成一步，前端对应页面会自动从「待实现」变成有数据，可以直接在浏览器里验收。

## 排查

| 现象 | 原因 / 处理 |
|------|-------------|
| 前端提示「后端服务未启动」 | 后端没跑，`cd backend && make run` |
| 前端接口报 501 | 该模块后端尚未实现，属预期；按 TODO 实现即可 |
| 端口被占用 | 改 `backend/.env` 的 `SERVER_PORT`，或 `lsof -i:8080` 查占用进程 |
| 写库报 `disk I/O error` / 连 `GET /categories` 也 500 | 多半是服务被**强杀**（IDE 直接停任务 / `kill -9`）后残留了热日志。先停掉所有 `blog-api` 进程，删掉 `backend/blog.db-journal`，再重启；仍不行就 `make reset-db` 重建（种子会自动重灌）。**切勿只删 `blog.db-journal` 而不重启进程** |
| CGO / SDK 编译报错 | 重装 CommandLineTools；仍不行可换纯 Go 驱动 `github.com/glebarez/sqlite`（注意 API 基本兼容） |
| 想重置数据 | `cd backend && make reset-db`，再重启后端 |

> 实测结论：在 `/Volumes/out_share`（APFS）上，**用编译好的二进制正常启动 / SIGTERM 停止**，连续重启 3 次写入均正常（201）。所以数据库放在项目目录没有问题；唯一会踩的坑是**用 IDE 后台任务强杀进程**。

## 注意

- 越权校验必须在 **Service 层**强制做（作者 / 管理员），前端隐藏按钮只是体验层。
- 多表写入（文章 + 标签关联）必须包在事务里。
- 状态 / 角色一律用 `backend/constant/` 下的常量，禁止裸字符串。
