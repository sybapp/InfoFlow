package logic

import (
	"context"

	"github.com/sybapp/infoflow/applications/user/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByPhoneLogic {
	return &FindByPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByPhoneLogic) FindByPhone(in *pb.FindByPhoneRequest) (*pb.FindByPhoneResponse, error) {
	user, err := l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		return nil, err
	}

	return &pb.FindByPhoneResponse{
		UserId:   user.Id,
		Username: user.Username,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	}, nil
}
