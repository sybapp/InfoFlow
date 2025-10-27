package model

import (
    "context"
    "fmt"

    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeRecordModel = (*customLikeRecordModel)(nil)

type (
    // LikeRecordModel is an interface to be customized, add more methods here,
    // and implement the added methods in customLikeRecordModel.
    LikeRecordModel interface {
        likeRecordModel
        FindByBizAndTarget(ctx context.Context, bizId string, targetId, userId int64) (*LikeRecord, error)
    }

    customLikeRecordModel struct {
        *defaultLikeRecordModel
    }
)

// NewLikeRecordModel returns a model for the database table.
func NewLikeRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeRecordModel {
    return &customLikeRecordModel{
        defaultLikeRecordModel: newLikeRecordModel(conn, c, opts...),
    }
}

func (m *customLikeRecordModel) FindByBizAndTarget(ctx context.Context, bizId string, targetId, userId int64) (*LikeRecord, error) {
    var record LikeRecord
    query := fmt.Sprintf("select %s from %s where biz_id = ? and target_id = ? and user_id = ? limit 1", likeRecordRows, m.table)
    err := m.QueryRowNoCacheCtx(ctx, &record, query, bizId, targetId, userId)
    if err != nil {
        if err == sqlx.ErrNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &record, nil
}
