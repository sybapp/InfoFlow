package svc

import (
	"github.com/sybapp/infoflow/applications/applet/internal/config"
	"github.com/sybapp/infoflow/applications/user/rpc/user"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRpc  user.User
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	userCli, err := zrpc.NewClient(c.UserRPC)
	if err != nil {
		return nil, err
	}
	u := user.NewUser(userCli)

	r, err := redis.NewRedis(c.BizRedis)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:   c,
		UserRpc:  u,
		BizRedis: r,
	}, nil
}
