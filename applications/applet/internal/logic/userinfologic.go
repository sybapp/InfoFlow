package logic

import (
	"context"
	"encoding/json"

	"github.com/sybapp/infoflow/applications/applet/internal/svc"
	"github.com/sybapp/infoflow/applications/applet/internal/types"
	"github.com/sybapp/infoflow/applications/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		logx.Errorf("UserInfo error: %v", err)
		return nil, err
	}

	userInfo, err := l.svcCtx.UserRpc.FindById(l.ctx, &user.FindByIdRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("UserInfo error: %v", err)
		return nil, err
	}

	return &types.UserInfoResponse{
		UserId:   userInfo.UserId,
		Username: userInfo.Username,
		Phone:    userInfo.Phone,
		Avatar:   userInfo.Avatar,
	}, nil
}
