package logic

import (
	"context"
	"strconv"
	"strings"

	"github.com/sybapp/infoflow/applications/applet/internal/code"
	"github.com/sybapp/infoflow/applications/applet/internal/svc"
	"github.com/sybapp/infoflow/applications/applet/internal/types"
	"github.com/sybapp/infoflow/applications/user/rpc/user"
	"github.com/sybapp/infoflow/pkg/encrypt"
	"github.com/sybapp/infoflow/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Phone = strings.TrimSpace(req.Phone)
	if _, err := strconv.Atoi(req.Phone); err != nil {
		logx.Errorf("Register req: %v error: %v", req, code.RegisterPhoneNotNumber)
		return nil, code.RegisterPhoneNotNumber
	}

	if req.Username == "" {
		logx.Errorf("Register req: %v error: %v", req, code.RegisterUsernameEmpty)
		return nil, code.RegisterUsernameEmpty
	}

	if req.Phone == "" {
		logx.Errorf("Register req: %v error: %v", req, code.RegisterPhoneEmpty)
		return nil, code.RegisterPhoneEmpty
	}
	
	req.Password = encrypt.EncPassword(req.Password)
	regRet, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Username,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("Register req: %v error: %v", req, err)
		return nil, err
	}

	token, err := jwt.BuildTokens(
		jwt.TokenOptions{
			AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
			AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
			Fields: map[string]any{
				"userid": regRet.UserId,
			},
		},
	)
	if err != nil {
		logx.Errorf("Register req: %v error: %v", req, err)
		return nil, err
	}

	return &types.RegisterResponse{
		UserId: regRet.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

