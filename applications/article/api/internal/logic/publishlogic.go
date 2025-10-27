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

    val := l.ctx.Value("userId")
    if val == nil {
        logx.Error("Publish: userId not found in context")
        return nil, code.UnauthorizedError
    }

    num, ok := val.(json.Number)
    if !ok {
        logx.Errorf("Publish: userId type assertion failed, got type: %T", val)
        return nil, code.InvalidTokenError
    }

    userId, err := num.Int64()
    if err != nil {
        logx.Errorf("Publish: convert userId to int64 error: %v", err)
        return nil, code.InvalidTokenError
    }

    articleId, err := l.svcCtx.ArticleRpc.Publish(l.ctx, &article.PublishRequest{
        UserId:      userId,
        Title:       req.Title,
        Content:     req.Content,
        Cover:       req.Cover,
        Description: req.Description,
    })
    if err != nil {
        logx.Errorf("Publish error: %v, title: %s", err, req.Title)
        return nil, err
    }

    return &types.PublishResponse{
        ArticleId: articleId.ArticleId,
    }, nil
}
