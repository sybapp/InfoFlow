package logic

import (
	"context"

	"github.com/sybapp/infoflow/applications/user/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByUsernameLogic {
	return &FindByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByUsernameLogic) FindByUsername(in *pb.FindByUsernameRequest) (*pb.FindByUsernameResponse, error) {
	user, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		return nil, err
	}

	return &pb.FindByUsernameResponse{
		UserId:   user.Id,
		Username: user.Username,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	}, nil
}
