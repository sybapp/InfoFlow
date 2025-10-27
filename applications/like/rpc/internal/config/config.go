package config

import (
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
    zrpc.RpcServerConf
    DataSource string
    CacheRedis cache.CacheConf
    KqPusherConf struct {
        Brokers []string
        Topic   string
    }
}
