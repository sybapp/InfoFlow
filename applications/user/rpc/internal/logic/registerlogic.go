package logic

import (
	"context"
	"time"

	"github.com/sybapp/infoflow/applications/user/rpc/internal/model"
	"github.com/sybapp/infoflow/applications/user/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/user/rpc/pb"
	"github.com/sybapp/infoflow/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := l.checkUserExist(in.Username, in.Phone)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: in.Username,
		Phone:    in.Phone,
		Avatar:   in.Avatar,
		Password: in.Password,
		Ctime:    time.Now(),
		Mtime:    time.Now(),
	})
	if err != nil {
		logx.Errorf("Register req: %v error: %v", in, err)
		return nil, xcode.New(xcode.DatabaseError, "注册失败")
	}

	userId, err := ret.LastInsertId()
	if err != nil {
		logx.Errorf("Register req: %v error: %v", in, err)
		return nil, xcode.New(xcode.DatabaseError, "注册失败")
	}

	return &pb.RegisterResponse{UserId: userId}, nil
}

func (l *RegisterLogic) checkUserExist(username, phone string) error {
	user, err := l.svcCtx.UserModel.FindByUsername(l.ctx, username)
	if err != nil {
		logx.Errorf("Register req: %v error: %v", username, err)
		return xcode.New(xcode.DatabaseError, "查询用户失败")
	}
	if user != nil {
		logx.Errorf("Register req: %v error: %v", username, err)
		return xcode.New(xcode.UserHasRegistered, "用户名已经注册")
	}

	user, err = l.svcCtx.UserModel.FindByPhone(l.ctx, phone)
	if err != nil {
		logx.Errorf("Register req: %v error: %v", phone, err)
		return xcode.New(xcode.DatabaseError, "查询用户失败")
	}
	if user != nil {
		logx.Errorf("Register req: %v error: %v", phone, err)
		return xcode.New(xcode.UserHasRegistered, "手机号已经注册")
	}
	return nil
}
