package config

import (
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/zrpc"
)

type AuthConfig struct {
    AccessSecret string
    AccessExpire int64
}

type OssConfig struct {
    Endpoint         string
    AccessKeyId      string
    AccessKeySecret  string
    BucketName       string
    ConnectTimeout   int64 `json:",optional"`
    ReadWriteTimeout int64 `json:",optional"`
}

type Config struct {
    rest.RestConf
    Auth       AuthConfig
    Oss        OssConfig
    ArticleRPC zrpc.RpcClientConf
}
