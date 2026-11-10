package repositories

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"blog-api/config"
	"blog-api/constant"
	"blog-api/models/entity"
)

// InitDB 打开 SQLite、迁移表结构并灌入幂等种子数据
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		// 学习期打开 SQL 日志，便于观察 GORM 生成的语句
		Logger: logger.Default.LogMode(logger.Info),
		// SQLite 下关闭外键约束迁移，避免自关联 ParentID=0 触发约束错误
		DisableForeignKeyConstraintWhenMigrating: true,
		// 把驱动层的唯一约束错误翻译为 gorm.ErrDuplicatedKey，便于上层识别为 409
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Tag{},
		&entity.Article{},
		&entity.Comment{},
	); err != nil {
		return nil, fmt.Errorf("迁移数据库失败: %w", err)
	}

	if err := seed(db); err != nil {
		return nil, fmt.Errorf("灌入种子数据失败: %w", err)
	}

	return db, nil
}

// seed 幂等灌入种子数据：先 Count 再插入
func seed(db *gorm.DB) error {
	users, err := seedUsers(db)
	if err != nil {
		return err
	}
	categories, err := seedCategories(db)
	if err != nil {
		return err
	}
	tags, err := seedTags(db)
	if err != nil {
		return err
	}
	return seedArticles(db, users, categories, tags)
}

// seedUsers 灌入默认用户：admin / lukas
func seedUsers(db *gorm.DB) (map[string]entity.User, error) {
	result := make(map[string]entity.User)

	var count int64
	if err := db.Model(&entity.User{}).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		var existing []entity.User
		if err := db.Find(&existing).Error; err != nil {
			return nil, err
		}
		for _, user := range existing {
			result[user.Username] = user
		}
		return result, nil
	}

	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	lukasHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	users := []entity.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(adminHash),
			Nickname: "管理员",
			Role:     constant.RoleAdmin,
			Status:   constant.UserStatusActive,
		},
		{
			Username: "lukas",
			Email:    "lukas@example.com",
			Password: string(lukasHash),
			Nickname: "Lukas",
			Role:     constant.RoleUser,
			Status:   constant.UserStatusActive,
		},
	}
	if err := db.Create(&users).Error; err != nil {
		return nil, err
	}
	for _, user := range users {
		result[user.Username] = user
	}
	return result, nil
}

// seedCategories 灌入默认分类：go / db / engineering
func seedCategories(db *gorm.DB) (map[string]entity.Category, error) {
	result := make(map[string]entity.Category)

	var count int64
	if err := db.Model(&entity.Category{}).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		var existing []entity.Category
		if err := db.Find(&existing).Error; err != nil {
			return nil, err
		}
		for _, category := range existing {
			result[category.Slug] = category
		}
		return result, nil
	}

	categories := []entity.Category{
		{Name: "Go 语言", Slug: "go", Description: "Go 语言相关", SortOrder: 0},
		{Name: "数据库", Slug: "db", Description: "数据库相关", SortOrder: 1},
		{Name: "工程实践", Slug: "engineering", Description: "工程实践相关", SortOrder: 2},
	}
	if err := db.Create(&categories).Error; err != nil {
		return nil, err
	}
	for _, category := range categories {
		result[category.Slug] = category
	}
	return result, nil
}

// seedTags 灌入默认标签：GORM / Gin / SQLite / 并发 / 架构
func seedTags(db *gorm.DB) (map[string]entity.Tag, error) {
	result := make(map[string]entity.Tag)

	var count int64
	if err := db.Model(&entity.Tag{}).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		var existing []entity.Tag
		if err := db.Find(&existing).Error; err != nil {
			return nil, err
		}
		for _, tag := range existing {
			result[tag.Name] = tag
		}
		return result, nil
	}

	tags := []entity.Tag{
		{Name: "GORM", Slug: "gorm"},
		{Name: "Gin", Slug: "gin"},
		{Name: "SQLite", Slug: "sqlite"},
		{Name: "并发", Slug: "concurrency"},
		{Name: "架构", Slug: "architecture"},
	}
	if err := db.Create(&tags).Error; err != nil {
		return nil, err
	}
	for _, tag := range tags {
		result[tag.Name] = tag
	}
	return result, nil
}

// seedArticles 灌入 4 篇示例文章（3 篇已发布 + 1 篇草稿），并建立标签多对多关联
func seedArticles(db *gorm.DB, users map[string]entity.User, categories map[string]entity.Category, tags map[string]entity.Tag) error {
	var count int64
	if err := db.Model(&entity.Article{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	admin := users["admin"]
	lukas := users["lukas"]
	goCategory := categories["go"]
	dbCategory := categories["db"]
	engineeringCategory := categories["engineering"]

	now := time.Now()
	publishedAt1 := now.AddDate(0, 0, -10)
	publishedAt2 := now.AddDate(0, 0, -5)
	publishedAt3 := now.AddDate(0, 0, -2)

	seeds := []struct {
		article  entity.Article
		tagNames []string
	}{
		{
			article: entity.Article{
				Title:      "GORM 关联与预加载实战",
				Slug:       "gorm-association-preload",
				Summary:    "整理 GORM 的关联查询与 Preload 用法，帮助你避开 N+1 陷阱。",
				CoverImage: "",
				Content: `# GORM 关联与预加载实战

本文整理 GORM 中关联查询与预加载的常见用法，帮助你避开 N+1 陷阱。

## 核心概念

- Belongs To：多对一，如文章属于一个分类
- Has Many：一对多，如分类拥有多篇文章
- Many To Many：多对多，如文章与标签

## 预加载示例

    // 预加载作者、分类、标签，SQL 数量恒定
    db.Preload("Author").Preload("Category").Preload("Tags").Find(&articles)

## 小结

优先使用 Preload，而不是在循环里逐条查询。
`,
				Status:      constant.ArticleStatusPublished,
				ViewCount:   42,
				AuthorID:    admin.ID,
				CategoryID:  &goCategory.ID,
				PublishedAt: &publishedAt1,
			},
			tagNames: []string{"GORM", "架构"},
		},
		{
			article: entity.Article{
				Title:      "SQLite 在个人项目中的落地实践",
				Slug:       "sqlite-in-personal-projects",
				Summary:    "轻量级数据库 SQLite 的优势、限制与常见使用姿势。",
				CoverImage: "",
				Content: `# SQLite 在个人项目中的落地实践

SQLite 无需独立服务进程，非常适合个人项目与本地开发。

## 优点

- 零配置，单文件存储
- 事务完整，支持大部分 SQL 标准
- 迁移到 MySQL / PostgreSQL 成本较低

## 注意事项

    // 开启 WAL 提升并发读写表现
    PRAGMA journal_mode = WAL;

## 小结

先用 SQLite 快速起步，后续再按需切换数据库。
`,
				Status:      constant.ArticleStatusPublished,
				ViewCount:   18,
				AuthorID:    admin.ID,
				CategoryID:  &dbCategory.ID,
				PublishedAt: &publishedAt2,
			},
			tagNames: []string{"SQLite", "GORM"},
		},
		{
			article: entity.Article{
				Title:      "Gin 统一响应与中间件设计",
				Slug:       "gin-response-and-middleware",
				Summary:    "设计统一响应体与常用中间件，让接口风格保持一致性。",
				CoverImage: "",
				Content: `# Gin 统一响应与中间件设计

统一响应体能显著降低前后端的沟通成本。

## 统一响应体

- code：业务状态码
- message：提示信息
- data：业务数据
- error：错误详情

## 中间件顺序

    // 日志 -> 恢复 -> 跨域
    r.Use(middleware.Logger(), middleware.Recovery(), middleware.CORS())

## 小结

中间件应职责单一，便于组合与复用。
`,
				Status:      constant.ArticleStatusPublished,
				ViewCount:   27,
				AuthorID:    admin.ID,
				CategoryID:  &engineeringCategory.ID,
				PublishedAt: &publishedAt3,
			},
			tagNames: []string{"Gin", "架构"},
		},
		{
			article: entity.Article{
				Title:      "Go 并发编程踩坑记录（草稿）",
				Slug:       "go-concurrency-pitfalls",
				Summary:    "记录 channel 与 goroutine 使用过程中的常见问题。",
				CoverImage: "",
				Content: `# Go 并发编程踩坑记录

（草稿，内容完善中）

## 待整理要点

- goroutine 泄漏的常见成因
- channel 关闭的注意事项
- sync 包的正确姿势

    // 示例占位
    var wg sync.WaitGroup
`,
				Status:     constant.ArticleStatusDraft,
				ViewCount:  0,
				AuthorID:   lukas.ID,
				CategoryID: &goCategory.ID,
			},
			tagNames: []string{"并发"},
		},
	}

	for i := range seeds {
		item := &seeds[i]
		// 先创建文章主体（不级联保存标签），避免对已存在标签做重复写入
		if err := db.Omit("Tags").Create(&item.article).Error; err != nil {
			return err
		}
		if len(item.tagNames) > 0 {
			tagList := make([]entity.Tag, 0, len(item.tagNames))
			for _, name := range item.tagNames {
				tagList = append(tagList, tags[name])
			}
			if err := db.Model(&item.article).Association("Tags").Replace(tagList); err != nil {
				return err
			}
		}
	}
	return nil
}
