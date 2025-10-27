package logic

import (
    "context"

    "github.com/sybapp/infoflow/applications/article/rpc/internal/svc"
    "github.com/sybapp/infoflow/applications/article/rpc/pb"

    "github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
    return &ArticleDeleteLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *ArticleDeleteLogic) ArticleDelete(in *pb.ArticleDeleteRequest) (*pb.ArticleDeleteResponse, error) {
    if in.ArticleId <= 0 {
        logx.Errorf("ArticleDelete: invalid article id: %d", in.ArticleId)
        return &pb.ArticleDeleteResponse{}, nil
    }

    article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
    if err != nil {
        logx.Errorf("ArticleDelete: find article error: %v, articleId: %d", err, in.ArticleId)
        return nil, err
    }

    if article == nil {
        logx.Infof("ArticleDelete: article not found, articleId: %d", in.ArticleId)
        return &pb.ArticleDeleteResponse{}, nil
    }

    if article.AuthorId != in.UserId {
        logx.Errorf("ArticleDelete: permission denied, articleId: %d, userId: %d, authorId: %d", 
            in.ArticleId, in.UserId, article.AuthorId)
        return &pb.ArticleDeleteResponse{}, nil
    }

    err = l.svcCtx.ArticleModel.Delete(l.ctx, in.ArticleId)
    if err != nil {
        logx.Errorf("ArticleDelete: delete article error: %v, articleId: %d", err, in.ArticleId)
        return nil, err
    }

    return &pb.ArticleDeleteResponse{}, nil
}
