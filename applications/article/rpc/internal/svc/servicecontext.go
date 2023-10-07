package svc

import (
	"github.com/sybapp/infoflow/applications/article/rpc/internal/config"
	"github.com/sybapp/infoflow/applications/article/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/singleflight"
)

type ServiceContext struct {
	Config            config.Config
	ArticleModel      model.ArticleModel
	BizRedis          *redis.Redis
	SingleFlightGroup singleflight.Group
}

func NewServiceContext(c config.Config) *ServiceContext {
	aritcleModel := model.NewArticleModel(sqlx.NewMysql(c.DataSource), c.CacheRedis)
	redis := redis.MustNewRedis(c.BizRedis)

	return &ServiceContext{
		Config:       c,
		ArticleModel: aritcleModel,
		BizRedis:     redis,
	}
}
