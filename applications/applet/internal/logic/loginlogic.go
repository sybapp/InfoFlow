package logic

import (
	"context"

	"github.com/sybapp/infoflow/applications/applet/internal/svc"
	"github.com/sybapp/infoflow/applications/applet/internal/types"
	"github.com/sybapp/infoflow/applications/user/rpc/user"
	"github.com/sybapp/infoflow/pkg/jwt"
	"github.com/sybapp/infoflow/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
		Username:   req.Username,
		Password:   req.Password,
		Phone:      req.Phone,
	})
	if err != nil {
		return nil, err
	}

	if rpcResp.UserId > 0 {
		token, err := jwt.BuildTokens(jwt.TokenOptions{
			AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
			AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
			Fields: map[string]any{
				"userId": rpcResp.UserId,
			},
		})
		if err != nil {
			return nil, err
		}

		resp = &types.LoginResponse{
			UserId: rpcResp.UserId,
			Token: types.Token{
				AccessToken:  token.AccessToken,
				AccessExpire: token.AccessExpire,
			},
		}
		return resp, nil
	} else {
		return nil, xcode.New(xcode.ParameterError, "用户名或密码错误")
	}
}
