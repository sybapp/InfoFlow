package logic

import (
	"context"
	"encoding/json"

	"github.com/sybapp/infoflow/applications/article/api/internal/code"
	"github.com/sybapp/infoflow/applications/article/api/internal/svc"
	"github.com/sybapp/infoflow/applications/article/api/internal/types"
	"github.com/sybapp/infoflow/applications/article/rpc/article"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	minContentLen = 10
)

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (resp *types.PublishResponse, err error) {
	if len(req.Title) == 0 {
		return nil, code.ArtitleTitleEmpty
	}
	if len(req.Content) < minContentLen {
		return nil, code.ArticleContentTooFewWords
	}
	if len(req.Cover) == 0 {
		return nil, code.ArticleCoverEmpty
	}

	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("Publish get userId error: %v", err)
		return nil, err
	}

	articleId, err := l.svcCtx.AritcleRpc.Publish(l.ctx, &article.PublishRequest{
		UserId:      userId,
		Title:       req.Title,
		Content:     req.Content,
		Cover:       req.Cover,
		Description: req.Description,
	})
	if err != nil {
		logx.Errorf("Publish req: %v error: %v", req, err)
		return nil, err
	}

	return &types.PublishResponse{
		ArticleId: articleId.ArticleId,
	}, nil
}
