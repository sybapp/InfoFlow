package logic

import (
    "context"
    "strings"

    "github.com/sybapp/infoflow/applications/user/rpc/internal/model"
    "github.com/sybapp/infoflow/applications/user/rpc/internal/svc"
    "github.com/sybapp/infoflow/applications/user/rpc/pb"
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

func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
    in.Username = strings.TrimSpace(in.Username)
    in.Phone = strings.TrimSpace(in.Phone)

    if in.Username != "" {
        return l.loginByUsername(in)
    } else if in.Phone != "" {
        return l.loginByPhone(in)
    }

    return nil, xcode.New(xcode.ParameterError, "用户名或手机号不能为空")
}

func (l *LoginLogic) loginByUsername(in *pb.LoginRequest) (*pb.LoginResponse, error) {
    user, err := l.svcCtx.UserModel.FindByUsername(l.ctx, in.Username)
    if err != nil {
        logx.Errorf("Login by username error: %v, username: %s", err, in.Username)
        return nil, xcode.New(xcode.DatabaseError, "登录失败")
    }

    return l.verifyAndLogin(user, in.Password, in.Username)
}

func (l *LoginLogic) loginByPhone(in *pb.LoginRequest) (*pb.LoginResponse, error) {
    user, err := l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
    if err != nil {
        logx.Errorf("Login by phone error: %v, phone: %s", err, in.Phone)
        return nil, xcode.New(xcode.DatabaseError, "登录失败")
    }

    return l.verifyAndLogin(user, in.Password, in.Phone)
}

func (l *LoginLogic) verifyAndLogin(user *model.User, password, identifier string) (*pb.LoginResponse, error) {
    if user == nil {
        logx.Infof("Login failed: user not found, identifier: %s", identifier)
        return nil, xcode.New(xcode.ParameterError, "用户不存在")
    }

    if !encrypt.VerifyPassword(user.Password, password) {
        logx.Infof("Login failed: password mismatch, identifier: %s", identifier)
        return nil, xcode.New(xcode.ParameterError, "密码错误")
    }

    return &pb.LoginResponse{
        UserId: user.Id,
    }, nil
}
