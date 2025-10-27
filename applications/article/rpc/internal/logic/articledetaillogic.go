package logic

import (
    "context"

    "github.com/sybapp/infoflow/applications/article/rpc/internal/svc"
    "github.com/sybapp/infoflow/applications/article/rpc/pb"

    "github.com/zeromicro/go-zero/core/logx"
)

type ArticleDetailLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
    return &ArticleDetailLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *ArticleDetailLogic) ArticleDetail(in *pb.ArticleDetailRequest) (*pb.ArticleDetailResponse, error) {
    if in.ArticleId <= 0 {
        logx.Errorf("ArticleDetail: invalid article id: %d", in.ArticleId)
        return nil, nil
    }

    article, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.ArticleId)
    if err != nil {
        logx.Errorf("ArticleDetail: find article error: %v, articleId: %d", err, in.ArticleId)
        return nil, err
    }

    if article == nil {
        logx.Infof("ArticleDetail: article not found, articleId: %d", in.ArticleId)
        return &pb.ArticleDetailResponse{}, nil
    }

    return &pb.ArticleDetailResponse{
        Article: &pb.ArticleItem{
            Id:           article.Id,
            Title:        article.Title,
            Content:      article.Content,
            Description:  article.Description,
            Cover:        article.Cover,
            CommentCount: article.CommentNum,
            LikeCount:    article.LikeNum,
            PublishTime:  article.PublishTime.Unix(),
        },
    }, nil
}
