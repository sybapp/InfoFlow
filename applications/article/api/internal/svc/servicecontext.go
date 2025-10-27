package svc

import (
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    "github.com/sybapp/infoflow/applications/article/api/internal/config"
    "github.com/sybapp/infoflow/applications/article/rpc/article"
    "github.com/sybapp/infoflow/pkg/interceptors"
    "github.com/zeromicro/go-zero/zrpc"
)

const (
    defaultOssConnectTimeout   = 1
    defaultOssReadWriteTimeout = 3
)

type ServiceContext struct {
    Config     config.Config
    OssClient  *oss.Client
    ArticleRpc article.Article
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
    if c.Oss.ConnectTimeout == 0 {
        c.Oss.ConnectTimeout = defaultOssConnectTimeout
    }
    if c.Oss.ReadWriteTimeout == 0 {
        c.Oss.ReadWriteTimeout = defaultOssReadWriteTimeout
    }

    ossClient, err := oss.New(c.Oss.Endpoint, c.Oss.AccessKeyId, c.Oss.AccessKeySecret,
        oss.Timeout(c.Oss.ConnectTimeout, c.Oss.ReadWriteTimeout))
    if err != nil {
        return nil, err
    }

    articleCli := zrpc.MustNewClient(
        c.ArticleRPC,
        zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()),
    )

    articleRpc := article.NewArticle(articleCli)

    return &ServiceContext{
        Config:     c,
        OssClient:  ossClient,
        ArticleRpc: articleRpc,
    }, nil
}
