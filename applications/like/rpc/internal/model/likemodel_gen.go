// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	likeFieldNames          = builder.RawFieldNames(&Like{})
	likeRows                = strings.Join(likeFieldNames, ",")
	likeRowsExpectAutoSet   = strings.Join(stringx.Remove(likeFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	likeRowsWithPlaceHolder = strings.Join(stringx.Remove(likeFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	likeModel interface {
		Insert(ctx context.Context, data *Like) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Like, error)
		Update(ctx context.Context, data *Like) error
		Delete(ctx context.Context, id int64) error
	}

	defaultLikeModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Like struct {
		Id         int64     `db:"id"`          // 主键ID
		BizId      string    `db:"biz_id"`      // 业务ID
		TargetId   int64     `db:"target_id"`   // 点赞目标id
		UserId     int64     `db:"user_id"`     // 用户ID
		Type       int64     `db:"type"`        // 类型 0:点赞 1:点踩
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 最后修改时间
	}
)

func newLikeModel(conn sqlx.SqlConn) *defaultLikeModel {
	return &defaultLikeModel{
		conn:  conn,
		table: "`like`",
	}
}

func (m *defaultLikeModel) withSession(session sqlx.Session) *defaultLikeModel {
	return &defaultLikeModel{
		conn:  sqlx.NewSqlConnFromSession(session),
		table: "`like`",
	}
}

func (m *defaultLikeModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultLikeModel) FindOne(ctx context.Context, id int64) (*Like, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", likeRows, m.table)
	var resp Like
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultLikeModel) Insert(ctx context.Context, data *Like) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, likeRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.BizId, data.TargetId, data.UserId, data.Type)
	return ret, err
}

func (m *defaultLikeModel) Update(ctx context.Context, data *Like) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, likeRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.BizId, data.TargetId, data.UserId, data.Type, data.Id)
	return err
}

func (m *defaultLikeModel) tableName() string {
	return m.table
}