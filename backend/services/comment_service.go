package services

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"blog-api/constant"
	"blog-api/models/dto"
	"blog-api/models/entity"
	"blog-api/models/vo"
	"blog-api/repositories"
	"blog-api/utils"
)

// CommentService 评论业务逻辑
type CommentService struct {
	repo        *repositories.CommentRepository
	articleRepo *repositories.ArticleRepository
}

// NewCommentService 创建评论服务
func NewCommentService(repo *repositories.CommentRepository, articleRepo *repositories.ArticleRepository) *CommentService {
	return &CommentService{repo: repo, articleRepo: articleRepo}
}

// ListByArticle 评论列表（顶层 + 内嵌 replies，内存组装两级评论树）
//   - 先分页取顶层评论，再用一次 `parent_id IN (...)` 批量取回全部回复并按父 ID 分发写回，
//     避免在循环里逐条查询造成 N+1
func (s *CommentService) ListByArticle(ctx context.Context, articleID uint, page utils.PageQuery) ([]vo.CommentVO, int64, error) {
	topLevel, count, err := s.repo.ListTopLevel(ctx, articleID, page)
	if err != nil {
		return nil, 0, fmt.Errorf("查询评论列表失败: %w", err)
	}

	res := make([]vo.CommentVO, 0, len(topLevel))
	// index 记录「顶层评论 ID -> res 下标」，便于把回复 O(1) 分发回对应节点
	index := make(map[uint]int, len(topLevel))
	parentIDs := make([]uint, 0, len(topLevel))
	for _, comm := range topLevel {
		index[comm.ID] = len(res)
		parentIDs = append(parentIDs, comm.ID)
		res = append(res, toCommentVO(comm))
	}

	// 一次 IN 查询取回当前页所有顶层评论的回复
	replies, err := s.repo.ListReplies(ctx, parentIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("查询评论回复失败: %w", err)
	}
	for _, reply := range replies {
		if i, ok := index[reply.ParentID]; ok {
			res[i].Replies = append(res[i].Replies, toCommentVO(reply))
		}
	}
	return res, count, nil
}

// Create 发表评论 / 回复：校验文章已发布、父评论同文章、两级归并规则
func (s *CommentService) Create(ctx context.Context, articleID uint, req dto.CreateCommentRequest, userID uint) (*vo.CommentVO, error) {
	// 1) 文章必须存在且已发布
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NotFound("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}
	if article.Status != constant.ArticleStatusPublished {
		return nil, utils.Forbidden("文章未发布，无法评论")
	}

	// 2) 父评论必须属于同一篇文章；若父评论本身是回复，则归并到其顶层父评论（始终两级）
	parentID := req.ParentID
	if parentID != 0 {
		parent, err := s.repo.GetByID(ctx, parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, utils.NotFound("父评论不存在")
			}
			return nil, fmt.Errorf("查询父评论失败: %w", err)
		}
		if parent.ArticleID != articleID {
			return nil, utils.BadRequest("父评论不属于该文章")
		}
		if parent.ParentID != 0 {
			parentID = parent.ParentID
		}
	}

	comment := &entity.Comment{
		Content:   req.Content,
		ParentID:  parentID,
		UserID:    userID,
		ArticleID: articleID,
	}
	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("创建评论失败: %w", err)
	}

	// 3) 回读一次以拿到 Preload(User) 与时间戳，组装成响应 VO
	created, err := s.repo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, fmt.Errorf("读取评论失败: %w", err)
	}
	result := toCommentVO(*created)
	return &result, nil
}

// Delete 删除评论：三方权限判定（评论作者 / 文章作者 / 管理员），连带软删除子回复
func (s *CommentService) Delete(ctx context.Context, commentID uint, userID uint, role constant.Role) error {
	comm, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("评论不存在")
		}
		return fmt.Errorf("查询评论失败: %w", err)
	}
	article, err := s.articleRepo.GetByID(ctx, comm.ArticleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound("文章不存在")
		}
		return fmt.Errorf("查询文章失败: %w", err)
	}

	// 三方权限：管理员放行；否则必须是评论作者或文章作者（三者都不满足才拒绝）
	if role != constant.RoleAdmin && comm.UserID != userID && article.AuthorID != userID {
		return utils.Forbidden("权限不足")
	}

	// 顶层评论连带软删其回复，整个过程在仓储层同一事务内完成
	if err := s.repo.DeleteWithReplies(ctx, commentID); err != nil {
		return fmt.Errorf("删除评论失败: %w", err)
	}
	return nil
}

// toCommentVO 实体转评论 VO（Replies 初始化为非 nil 空切片，避免 JSON 出现 null）
func toCommentVO(comment entity.Comment) vo.CommentVO {
	return vo.CommentVO{
		ID:        comment.ID,
		Content:   comment.Content,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
		User: vo.UserBriefVO{
			ID:       comment.UserID,
			Username: comment.User.Username,
			Nickname: comment.User.Nickname,
			Avatar:   comment.User.Avatar,
		},
		Replies: make([]vo.CommentVO, 0),
	}
}
