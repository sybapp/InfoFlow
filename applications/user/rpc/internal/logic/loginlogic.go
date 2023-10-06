package logic

import (
	"context"
	"strings"

	"github.com/sybapp/infoflow/applications/user/rpc/internal/svc"
	"github.com/sybapp/infoflow/applications/user/rpc/service"
	"github.com/sybapp/infoflow/pkg/encrypt"
	"github.com/sybapp/infoflow/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *service.LoginRequest) (*service.LoginResponse, error) {
	in.Username = strings.TrimSpace(in.Username)
	in.Phone = strings.TrimSpace(in.Phone)

	if in.Username != "" {
		return l.loginByUsername(in)
	} else if in.Phone != "" {
		return l.loginByPhone(in)
	}

	return nil, xcode.New(xcode.ParameterError, "用户名或手机号不能为空")
}

func (l *LoginLogic) loginByUsername(in *service.LoginRequest) (*service.LoginResponse, error) {
	user, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
	if err != nil {
		logx.Errorf("Login req: %v error: %v", in, err)
		return nil, xcode.New(xcode.DatabaseError, "登录失败")
	}

	if user == nil {
		logx.Errorf("Login req: %v error: %v", in, "用户不存在")
		return nil, xcode.New(xcode.ParameterError, "用户不存在")
	}

	if user.Password != encrypt.EncPassword(in.Password) {
		logx.Errorf("Login req: %v error: %v", in, "密码错误")
		return nil, xcode.New(xcode.ParameterError, "密码错误")
	}

	return &service.LoginResponse{
		UserId: user.Id,
	}, nil
}

func (l *LoginLogic) loginByPhone(in *service.LoginRequest) (*service.LoginResponse, error) {
	user, err := l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		logx.Errorf("Login req: %v error: %v", in, err)
		return nil, xcode.New(xcode.DatabaseError, "登录失败")
	}

	if user == nil {
		logx.Errorf("Login req: %v error: %v", in, "用户不存在")
		return nil, xcode.New(xcode.ParameterError, "用户不存在")
	}

	if user.Password != encrypt.EncPassword(in.Password) {
		logx.Errorf("Login req: %v error: %v", in, "密码错误")
		return nil, xcode.New(xcode.ParameterError, "密码错误")
	}

	return &service.LoginResponse{
		UserId: user.Id,
	}, nil
}
