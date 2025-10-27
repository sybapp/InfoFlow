package logic

import (
    "context"
    "encoding/json"

    "github.com/sybapp/infoflow/applications/like/rpc/internal/svc"
    "github.com/sybapp/infoflow/applications/like/rpc/internal/types"
    "github.com/sybapp/infoflow/applications/like/rpc/service"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/threading"
)

type ThumbupLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
    return &ThumbupLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *ThumbupLogic) Thumbup(in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
    if in.ObjId <= 0 || in.UserId <= 0 || in.BizId == "" {
        logx.Errorf("Thumbup: invalid params, bizId: %s, objId: %d, userId: %d", 
            in.BizId, in.ObjId, in.UserId)
        return &service.ThumbupResponse{}, nil
    }

    record, err := l.svcCtx.LikeRecordModel.FindByBizAndTarget(l.ctx, in.BizId, in.ObjId, in.UserId)
    if err != nil {
        logx.Errorf("Thumbup: query record error: %v, bizId: %s, objId: %d, userId: %d", 
            err, in.BizId, in.ObjId, in.UserId)
        return nil, err
    }

    if record != nil && record.Type == int64(in.LikeType) {
        logx.Infof("Thumbup: already liked, bizId: %s, objId: %d, userId: %d, type: %d",
            in.BizId, in.ObjId, in.UserId, in.LikeType)
        return &service.ThumbupResponse{}, nil
    }

    msg := &types.ThumbupMsg{
        BizId:    in.BizId,
        ObjId:    in.ObjId,
        UserId:   in.UserId,
        LikeType: in.LikeType,
    }

    data, err := json.Marshal(msg)
    if err != nil {
        logx.Errorf("Thumbup: marshal msg error: %v, msg: %+v", err, msg)
        return nil, err
    }

    threading.GoSafe(func() {
        maxRetries := 3
        for i := 0; i < maxRetries; i++ {
            err := l.svcCtx.KqPusherClient.Push(string(data))
            if err == nil {
                logx.Infof("Thumbup: successfully pushed to kafka, bizId: %s, objId: %d, userId: %d",
                    in.BizId, in.ObjId, in.UserId)
                return
            }
            logx.Errorf("Thumbup: kq push error (attempt %d/%d): %v, data: %s", 
                i+1, maxRetries, err, string(data))
        }
        logx.Errorf("Thumbup: failed to push to kafka after %d retries, data: %s", 
            maxRetries, string(data))
    })

    return &service.ThumbupResponse{}, nil
}
