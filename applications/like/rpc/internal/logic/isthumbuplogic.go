package logic

import (
    "context"

    "github.com/sybapp/infoflow/applications/like/rpc/internal/svc"
    "github.com/sybapp/infoflow/applications/like/rpc/service"

    "github.com/zeromicro/go-zero/core/logx"
)

type IsThumbupLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewIsThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsThumbupLogic {
    return &IsThumbupLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *IsThumbupLogic) IsThumbup(in *service.IsThumbupRequest) (*service.IsThumbupResponse, error) {
    if in.TargetId <= 0 || in.UserId <= 0 || in.BizId == "" {
        logx.Errorf("IsThumbup: invalid params, bizId: %s, targetId: %d, userId: %d", 
            in.BizId, in.TargetId, in.UserId)
        return &service.IsThumbupResponse{}, nil
    }

    record, err := l.svcCtx.LikeRecordModel.FindByBizAndTarget(l.ctx, in.BizId, in.TargetId, in.UserId)
    if err != nil {
        logx.Errorf("IsThumbup: query error: %v, bizId: %s, targetId: %d, userId: %d", 
            err, in.BizId, in.TargetId, in.UserId)
        return nil, err
    }

    userThumbup := &service.UserThumbup{}
    if record != nil {
        userThumbup.UserId = record.UserId
        userThumbup.LikeType = int32(record.Type)
        userThumbup.ThumbupTime = record.CreateTime.Unix()
    }

    return &service.IsThumbupResponse{
        UserThumbups: map[int64]*service.UserThumbup{
            in.TargetId: userThumbup,
        },
    }, nil
}
