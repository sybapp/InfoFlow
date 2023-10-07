package logic

import (
	"context"
	"time"

	"github.com/sybapp/infoflow/applications/article/rpc/internal/model"
	"github.com/sybapp/infoflow/applications/article/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *pb.PublishRequest) (*pb.PublishResponse, error) {
	ret, err := l.svcCtx.AritcleModel.Insert(l.ctx, &model.Article{
		AuthorId:    in.UserId,
		Title:       in.Title,
		Content:     in.Content,
		Cover:       in.Cover,
		Description: in.Description,
		PublishTime: time.Now(),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	})
	if err != nil {
		logx.Errorf("Publish insert req: %v error: %v", in, err)
		return nil, err
	}

	articleId, err := ret.LastInsertId()
	if err != nil {
		logx.Errorf("Publish lastInsertId req: %v error: %v", in, err)
		return nil, err
	}

	return &pb.PublishResponse{ArticleId: articleId}, nil
}
