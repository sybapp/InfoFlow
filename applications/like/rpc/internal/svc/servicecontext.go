package svc

import (
    "github.com/sybapp/infoflow/applications/like/rpc/internal/config"
    "github.com/sybapp/infoflow/applications/like/rpc/internal/model"

    "github.com/zeromicro/go-queue/kq"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config          config.Config
    KqPusherClient  *kq.Pusher
    LikeRecordModel model.LikeRecordModel
    LikeCountModel  model.LikeCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.DataSource)
    return &ServiceContext{
        Config:          c,
        KqPusherClient:  kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
        LikeRecordModel: model.NewLikeRecordModel(conn, c.CacheRedis),
        LikeCountModel:  model.NewLikeCountModel(conn, c.CacheRedis),
    }
}
