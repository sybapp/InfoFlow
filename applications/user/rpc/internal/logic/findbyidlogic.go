package logic

import (
	"context"

	"github.com/sybapp/infoflow/applications/user/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *service.FindByIdRequest) (*service.FindByIdResponse, error) {
	user, err := l.svcCtx.UserModel.FindById(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return &service.FindByIdResponse{
		UserId:   user.Id,
		Username: user.Username,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	}, nil
}
