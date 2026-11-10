package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"blog-api/config"
	"blog-api/handlers"
	"blog-api/middleware"
	"blog-api/repositories"
	"blog-api/services"
)

// Setup 组装依赖并注册全部路由
func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	// 仓储层
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	articleRepo := repositories.NewArticleRepository(db)
	commentRepo := repositories.NewCommentRepository(db)
	statsRepo := repositories.NewStatsRepository(db)

	// 业务层
	authService := services.NewAuthService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	tagService := services.NewTagService(tagRepo)
	articleService := services.NewArticleService(articleRepo)
	commentService := services.NewCommentService(commentRepo, articleRepo)
	statsService := services.NewStatsService(statsRepo)

	// 处理层
	authHandler := handlers.NewAuthHandler(authService, cfg)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	tagHandler := handlers.NewTagHandler(tagService)
	articleHandler := handlers.NewArticleHandler(articleService)
	commentHandler := handlers.NewCommentHandler(commentService)
	statsHandler := handlers.NewStatsHandler(statsService)

	// 中间件
	authMiddleware := middleware.Auth([]byte(cfg.JWTSecret))
	optionalAuth := middleware.OptionalAuth([]byte(cfg.JWTSecret))
	adminOnly := middleware.AdminOnly()
	// 注意：管理员路由必须同时挂 authMiddleware + adminOnly，
	// AdminOnly 依赖 Auth 写入上下文的角色，只挂 adminOnly 会永远 403

	engine := gin.New()
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CORS())

	// 健康检查（顶层，不在 /api/v1 下）
	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	api := engine.Group("/api/v1")

	// 认证模块
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", authMiddleware, authHandler.Me)
		auth.PUT("/me", authMiddleware, authHandler.UpdateMe)
		auth.PUT("/password", authMiddleware, authHandler.ChangePassword)
	}

	// 分类模块：列表公开，写操作仅管理员
	categories := api.Group("/categories")
	{
		categories.GET("", categoryHandler.List)
		categories.POST("", authMiddleware, adminOnly, categoryHandler.Create)
		categories.PUT("/:id", authMiddleware, adminOnly, categoryHandler.Update)
		categories.DELETE("/:id", authMiddleware, adminOnly, categoryHandler.Delete)
	}

	// 标签模块：列表公开，写操作仅管理员
	tags := api.Group("/tags")
	{
		tags.GET("", tagHandler.List)
		tags.POST("", authMiddleware, adminOnly, tagHandler.Create)
		tags.PUT("/:id", authMiddleware, adminOnly, tagHandler.Update)
		tags.DELETE("/:id", authMiddleware, adminOnly, tagHandler.Delete)
	}

	// 文章模块：读接口可选鉴权（游客可见已发布，作者/管理员可见草稿），写操作需登录
	articles := api.Group("/articles")
	{
		// 注意：archive / hot 必须注册在 :id 之前，避免被 :id 抢占
		articles.GET("/archive", optionalAuth, articleHandler.Archive)
		articles.GET("/hot", optionalAuth, articleHandler.Hot)
		articles.GET("", optionalAuth, articleHandler.List)
		articles.GET("/:id", optionalAuth, articleHandler.GetByID)
		articles.GET("/:id/comments", optionalAuth, commentHandler.ListByArticle)
		articles.POST("", authMiddleware, articleHandler.Create)
		articles.POST("/:id/comments", authMiddleware, commentHandler.Create)
		articles.PUT("/:id", authMiddleware, articleHandler.Update)
		articles.PATCH("/:id/status", authMiddleware, articleHandler.ChangeStatus)
		articles.DELETE("/:id", authMiddleware, articleHandler.Delete)
	}

	// 评论删除
	comments := api.Group("/comments")
	{
		comments.DELETE("/:id", authMiddleware, commentHandler.Delete)
	}

	// 统计模块：仅管理员
	stats := api.Group("/stats")
	stats.Use(authMiddleware, adminOnly)
	{
		stats.GET("/overview", statsHandler.Overview)
		stats.GET("/categories", statsHandler.Categories)
	}

	return engine
}
