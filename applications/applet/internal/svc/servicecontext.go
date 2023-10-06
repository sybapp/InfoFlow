package svc

import (
	"github.com/sybapp/infoflow/applications/applet/internal/config"
	"github.com/sybapp/infoflow/applications/user/rpc/user"
	"github.com/sybapp/infoflow/pkg/interceptors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRpc  user.User
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	userCli := zrpc.MustNewClient(
		c.UserRPC,
		zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()),
	)

	u := user.NewUser(userCli)
	r := redis.MustNewRedis(c.BizRedis)

	return &ServiceContext{
		Config:   c,
		UserRpc:  u,
		BizRedis: r,
	}, nil
}
